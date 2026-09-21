import { useEffect, useState } from "react";
import VideoPlayer from "./VideoPlayer";
type Video = { id: string; name: string; status: string; url?: string; error?: string };
const api = "http://localhost:8080";
export default function App() {
  const [videos, setVideos] = useState<Video[]>([]);
  const [error, setError] = useState("");
  const [uploading, setUploading] = useState(false);
  useEffect(() => {
    const controller = new AbortController();
    const refresh = async () => {
      try { const r = await fetch(`${api}/api/videos`, { signal: controller.signal }); if (!r.ok) throw new Error("Could not load videos"); setVideos(await r.json()); }
      catch (e) { if (!controller.signal.aborted) setError(String(e)); }
    };
    void refresh(); const timer = setInterval(refresh, 1500);
    return () => { controller.abort(); clearInterval(timer); };
  }, []);
  return <main style={{ maxWidth: 900, margin: "2rem auto", fontFamily: "sans-serif" }}>
    <h1>HLS video library</h1><p>Upload a video up to 100 MiB. One video processes at a time.</p>
    <input type="file" accept="video/*" disabled={uploading} onChange={async e => {
      const file = e.target.files?.[0]; if (!file) return; setUploading(true); setError("");
      try { const body = new FormData(); body.append("file", file); const r = await fetch(`${api}/api/videos`, { method: "POST", body }); if (!r.ok) throw new Error(await r.text()); const v: Video = await r.json(); setVideos(old => [...old.filter(x => x.id !== v.id), v]); }
      catch (e) { setError(String(e)); } finally { setUploading(false); }
    }} />{uploading && <p>Uploading…</p>}{error && <p role="alert">{error}</p>}
    {videos.map(v => <section key={v.id}><h2>{v.name}</h2><p>{v.status}{v.error ? `: ${v.error}` : ""}</p>{v.status === "completed" && v.url && <VideoPlayer src={api + v.url} />}</section>)}
  </main>;
}
