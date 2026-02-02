import { AddressInfo } from "node:net"
import type{ Plugin} from "vite"
type HoshikuzuDevtoolsOptions={
  /**
   * 通知地址
   */
  notifyUrl:string
}
export function hoshikuzuDevtools(options:HoshikuzuDevtoolsOptions):Plugin{
  return {
    name:"vite-plugin-hoshikuzu-devtools",
    configureServer(server){
      server.httpServer?.on('listening',async()=>{
        try{
          const addr = server.httpServer!.address() as AddressInfo
          const host = addr.family === 'IPv6'
              ? '[::1]'
              : (addr.address === '0.0.0.0' ? 'localhost' : addr.address);
              const port = addr.port
              const httpUrl = `http://${host}:${port}`
              const wsUrl = `ws://${host}:${port}`

              const payload = {
                http: httpUrl,
                ws: wsUrl,
                host: host,
                port: port,
                time_unix: Date.now(),
              };

              await fetch(options.notifyUrl+"/online", {
              method: 'POST',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify(payload, null, 2),
            });
        }catch(err){
          console.warn(err)
        }
      })
    },
    handleHotUpdate(ctx) {
        const { file, modules } = ctx;
        // 只上报你关心的源码类型
        const allowExts = ['.ts', '.js', '.tsx', '.jsx', '.vue', '.scss', '.css'];
        const isTargetFile = allowExts.some(ext => file.endsWith(ext));

        if (isTargetFile) {
          const payload = {
            files: [file], // 可收集多个变更文件
            time: Date.now(),
          };

          // 异步上报，不阻塞 HMR
          fetch(options.notifyUrl+"/update", {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
          }).catch(err => {
            console.warn(err);
          });
        }

        return modules;
      },
  }
}
