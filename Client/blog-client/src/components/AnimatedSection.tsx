import { type ReactNode } from "react";
import Box from "@mui/material/Box";
import Slide from "@mui/material/Slide";
import { useInView } from "react-intersection-observer";

type Direction = "left" | "right";

interface Props {
  children: ReactNode;
  direction: Direction;
  id?: string;
}

export default function AnimatedSection({ children, direction, id }: Props) {
  const { ref, inView } = useInView({
    threshold: 0.3,
  });

  return (
    <Box
      id={id}
      ref={ref}
      sx={{
        minHeight: "60vh",
        display: "flex",
        alignItems: "center",
      }}
    >
      <Slide
        direction={direction}
        in={inView}
        timeout={700}
        mountOnEnter
        unmountOnExit
      >
        <Box width="100%">{children}</Box>
      </Slide>
    </Box>
  );
}
