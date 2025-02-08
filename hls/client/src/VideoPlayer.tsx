import { useEffect, useRef } from "react";
import videojs from "video.js";
import "video.js/dist/video-js.css";

const VideoPlayer: React.FC = () => {
  const videoRef = useRef<HTMLVideoElement | null>(null);

  const playerRef = useRef(null);

  useEffect(() => {
    if (!videoRef.current) return;
    //@ts-ignore
    playerRef.current = videojs(videoRef.current, {
      controls: true,
      autoplay: false,
      preload: "auto",
      fluid: true, // Makes the player responsive
    });

    return () => {
      if (playerRef.current) {
        //@ts-ignore
        playerRef.current.dispose();
      }
    };
  }, []);

  return (
    <div data-vjs-player>
      <video ref={videoRef} className="video-js vjs-default-skin">
        <source
          src="http://localhost:8080/hls/output.m3u8"
          type="application/x-mpegURL"
        />
        Your browser does not support the video tag.
      </video>
    </div>
  );
};

export default VideoPlayer;
