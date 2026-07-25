const http = require('http');
const fs = require('fs');
const path = require('path');
const API_TARGET = 'http://127.0.0.1:11434';

const server = http.createServer(async (req,res)=>{
  if (req.url === '/' || req.url.endsWith('.html') || req.url.endsWith('.js') || req.url.endsWith('.css')) {
    const p = path.join(__dirname, req.url === '/' ? 'ollama-ui.html' : req.url);
    if (!fs.existsSync(p)) { res.writeHead(404); res.end('not found'); return; }
    res.writeHead(200, {'Content-Type': req.url.endsWith('.html') ? 'text/html' : 'application/octet-stream'});
    fs.createReadStream(p).pipe(res); return;
  }
  // Proxy other requests to Ollama with CORS
  const upstream = API_TARGET + req.url;
  const uReq = http.request(upstream, { method: req.method, headers: req.headers }, uRes=>{
    res.writeHead(uRes.statusCode, { ...uRes.headers,
      'access-control-allow-origin':'*',
      'access-control-allow-methods':'GET,POST,OPTIONS',
      'access-control-allow-headers':'Content-Type' });
    uRes.pipe(res);
  });
  req.pipe(uReq);
  uReq.on('error', e=>{ res.writeHead(502); res.end('bad gateway'); });
});
server.listen(8000, ()=>console.log('listening 8000'));
