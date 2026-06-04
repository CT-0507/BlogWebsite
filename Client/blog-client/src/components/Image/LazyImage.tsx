import type { BoxProps } from "@mui/material";
import Box from "@mui/material/Box";
import { useEffect, useRef, useState } from "react";

// interface LazyImageProps {
//   sx?: SxProps<Theme>;
//   src?: string;
//   alt?: string;
//   placeholder?: string;
//   className?: string;
//   id?: string;
// }

interface LazyImageProps extends BoxProps<"img"> {
  placeholder?: string;
}

export default function LazyImage({ ...params }: LazyImageProps) {
  const imgRef = useRef<HTMLImageElement | null>(null);
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    if (isVisible) return;

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsVisible(true);
          observer.unobserve(entry.target);
        }
      },
      {
        threshold: 0.1,
        rootMargin: "100px",
      },
    );

    const currentRef = imgRef.current;

    if (currentRef) {
      observer.observe(currentRef);
    }

    return () => {
      if (currentRef) {
        observer.unobserve(currentRef);
      }
    };
  }, [isVisible]);

  return <Box component="img" ref={imgRef} loading="lazy" {...params} />;
}
