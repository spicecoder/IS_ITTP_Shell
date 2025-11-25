// intention-client.js
const crypto = require('crypto');
const axios = require('axios');

function sha256hex(input){
  return crypto.createHash('sha256').update(input).digest('hex');
}
function hmacSign(msg, secret){
  return crypto.createHmac('sha256', secret).update(msg).digest('hex');
}

module.exports = function createClient({ serverUrl, sharedSecret }){
  if(!serverUrl) throw new Error('serverUrl required');
  if(!sharedSecret) throw new Error('sharedSecret required (demo only)');

  async function sendIntention(intentName, payload, opts = {}){
    const canonicalPayload = JSON.stringify(payload);
    const contentHash = 'sha256:' + sha256hex(canonicalPayload);
    const intention = {
      intent: intentName,
      contentHash,
      clientId: opts.clientId || 'client:unknown',
      createdAt: (new Date()).toISOString()
    };
    const canonicalIntention = JSON.stringify(intention);
    const signature = hmacSign(canonicalIntention, sharedSecret);
    intention.signature = signature;

    const resp = await axios.post(`${serverUrl}/intention/submit`, { intention, payload }, { timeout: 10000 });
    return resp.data;
  }

  return { sendIntention };
};
