// test-send.js
const clientFactory = require('./intention-client');
const client = clientFactory({ serverUrl:'http://localhost:3000', sharedSecret:'topsecret' });

(async ()=>{
  try {
    console.log('Sending innocent post...');
    console.log(await client.sendIntention('publish_post', { text: 'this is innocent learning material' }, { clientId:'user:alice' }));
  } catch(e){
    console.error('Error sending innocent post', e.response ? e.response.data : e.message);
  }

  try {
    console.log('Sending flagged post...');
    console.log(await client.sendIntention('publish_post', { text: 'bomb plan and kill' }, { clientId:'user:attacker' }));
  } catch(e){
    console.error('Error sending flagged post', e.response ? e.response.data : e.message);
  }
})();
