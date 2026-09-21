import { useEffect, useRef } from "react";
import videojs from "video.js";
import "video.js/dist/video-js.css";
export default function VideoPlayer({ src }: { src: string }) {
  const container = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (!container.current) return;
    const element = document.createElement("video-js");
    container.current.appendChild(element);
    const player = videojs(element, { controls: true, fluid: true, preload: "metadata", sources: [{ src, type: "application/x-mpegURL" }] });
    return () => { player.dispose(); };
  }, [src]);
  return <div data-vjs-player ref={container} />;
}
