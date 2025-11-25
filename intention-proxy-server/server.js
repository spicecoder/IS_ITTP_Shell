// server.js
// Demo Intention Space Object proxy (see README)
const express = require('express');
const bodyParser = require('body-parser');
const crypto = require('crypto');
const fs = require('fs');
const path = require('path');

const PORT = process.env.PORT || 3000;
const SHARED_SECRET = process.env.SHARED_SECRET || 'dev_secret';
const TRACE_FILE = path.join(__dirname, 'cpuX_traces.jsonl');
const PUBLISH_DIR = path.join(__dirname, 'published');
if(!fs.existsSync(PUBLISH_DIR)) fs.mkdirSync(PUBLISH_DIR);

const app = express();
app.use(bodyParser.json({limit:'5mb'}));

let quarantine = [];
let published = [];

function fastClassifier(payload){
  const text = JSON.stringify(payload).toLowerCase();
  const flags = ['kill', 'bomb', 'suicide', 'groom', 'hate', 'pedo', 'childporn', 'terror'];
  let score = 0;
  for(const f of flags){
    if(text.includes(f)) score += 0.3;
  }
  return Math.min(1, score);
}

function sha256hex(input){
  return crypto.createHash('sha256').update(input).digest('hex');
}
function hmacSign(msg, secret){
  return crypto.createHmac('sha256', secret).update(msg).digest('hex');
}
function writeTrace(trace){
  fs.appendFileSync(TRACE_FILE, JSON.stringify(trace) + '\n');
}

function verifyIntentionMiddleware(req, res, next){
  const intention = req.body.intention;
  const payload = req.body.payload;
  if(!intention || !payload) return res.status(400).json({error:'missing intention or payload'});

  const actualHash = sha256hex(JSON.stringify(payload));
  if(!intention.contentHash || intention.contentHash.replace(/^sha256:/,'') !== actualHash){
    const trace = {time: new Date().toISOString(), event:'hash_mismatch', intention, ip:req.ip};
    writeTrace(trace);
    return res.status(400).json({error:'contentHash mismatch'});
  }

  const { signature } = intention;
  if(!signature) return res.status(400).json({error:'missing signature'});
  const copy = Object.assign({}, intention);
  delete copy.signature;
  const canonical = JSON.stringify(copy);
  const expected = hmacSign(canonical, SHARED_SECRET);
  try{
    if(!crypto.timingSafeEqual(Buffer.from(expected,'hex'), Buffer.from(signature,'hex'))){
      const trace = {time: new Date().toISOString(), event:'signature_mismatch', intention, ip:req.ip};
      writeTrace(trace);
      return res.status(401).json({error:'signature verification failed'});
    }
  } catch(e){
    // timingSafeEqual throws if lengths mismatch
    const trace = {time: new Date().toISOString(), event:'signature_mismatch_exception', intention, ip:req.ip};
    writeTrace(trace);
    return res.status(401).json({error:'signature verification failed'});
  }

  req.verified = {intention, payload};
  next();
}

app.post('/intention/submit', verifyIntentionMiddleware, (req,res)=>{
  const { intention, payload } = req.verified;
  const score = fastClassifier(payload);

  const trace = {
    time: new Date().toISOString(),
    event: 'intention_received',
    intentName: intention.intent,
    clientId: intention.clientId || null,
    contentHash: intention.contentHash,
    classifierScore: score,
  };

  const QUARANTINE_THRESHOLD = 0.5;
  if(score >= QUARANTINE_THRESHOLD){
    const qid = crypto.randomBytes(6).toString('hex');
    const entry = { qid, intention, payload, receivedAt: new Date().toISOString(), classifierScore: score };
    quarantine.push(entry);
    trace.event = 'quarantined';
    trace.qid = qid;
    writeTrace(trace);
    return res.status(202).json({status:'quarantined', qid});
  }

  const pubId = crypto.randomBytes(6).toString('hex');
  const filename = path.join(PUBLISH_DIR, `${pubId}.json`);
  fs.writeFileSync(filename, JSON.stringify({intention,payload,publishedAt:new Date().toISOString()}, null, 2));
  published.push({pubId, filename, intention});
  trace.event = 'published';
  trace.pubId = pubId;
  writeTrace(trace);
  return res.status(200).json({status:'published', pubId});
});

app.get('/moderation/pending', (req,res)=>{
  const list = quarantine.map(q => ({ qid: q.qid, intent: q.intention.intent, client: q.intention.clientId, receivedAt: q.receivedAt, score: q.classifierScore }));
  res.json({pending: list});
});

app.get('/moderation/item/:qid', (req,res)=>{
  const q = quarantine.find(x => x.qid === req.params.qid);
  if(!q) return res.status(404).json({error:'not found'});
  res.json({intention:q.intention, payload:q.payload, receivedAt:q.receivedAt, score:q.classifierScore});
});

app.post('/moderation/decide', bodyParser.json(), (req,res)=>{
  const { qid, decision, moderator, note } = req.body;
  const idx = quarantine.findIndex(x => x.qid === qid);
  if(idx === -1) return res.status(404).json({error:'qid not found'});
  const item = quarantine.splice(idx,1)[0];
  const trace = { time: new Date().toISOString(), event:'moderation_decision', qid, decision, moderator, note };
  if(decision === 'publish'){
    const pubId = crypto.randomBytes(6).toString('hex');
    const filename = path.join(PUBLISH_DIR, `${pubId}.json`);
    fs.writeFileSync(filename, JSON.stringify({intention:item.intention,payload:item.payload,publishedAt:new Date().toISOString(),moderator}, null, 2));
    published.push({pubId, filename, intention:item.intention});
    trace.pubId = pubId;
  } else {
    trace.rejected = true;
  }
  writeTrace(trace);
  res.json({status:'ok', trace});
});

app.get('/published', (req,res) => res.json({published: published.map(p => ({pubId:p.pubId, filename:p.filename, intent:p.intention.intent}))}));
app.get('/health', (req,res) => res.json({ok:true, pending: quarantine.length}));

app.listen(PORT, ()=> console.log(`Intention Proxy listening on ${PORT}`));
