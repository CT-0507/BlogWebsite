import SpeedDial from "@mui/material/SpeedDial";
import SpeedDialAction from "@mui/material/SpeedDialAction";
import SpeedDialIcon from "@mui/material/SpeedDialIcon";
import ShareIcon from "@mui/icons-material/Share";
import FacebookIcon from "@mui/icons-material/Facebook";
import CheckIcon from "@mui/icons-material/Check";
import { useState } from "react";

// const actions = [
//   {
//     icon: <FacebookIcon />,
//     name: "Facebook",
//     onClick: async () => {
//       //   window.open(
//       //     `https://www.facebook.com/sharer/sharer.php?u=${window.location.href}`
//       //   );
//       try {
//         await navigator.clipboard.writeText(window.location.href);
//       } catch (err) {
//         console.error("Copy failed", err);
//       }
//     },
//   },
//   {
//     icon: <XIcon />,
//     name: "Twitter",
//     onClick: () =>
//       window.open(
//         `https://twitter.com/intent/tweet?url=${window.location.href}`
//       ),
//   },
//   {
//     icon: <LinkedInIcon />,
//     name: "LinkedIn",
//     onClick: () =>
//       window.open(
//         `https://www.linkedin.com/sharing/share-offsite/?url=${window.location.href}`
//       ),
//   },
// ];

export default function SocialShareDial() {
  const [isCopied, setIsCopy] = useState(false);
  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(window.location.href);
      setIsCopy(true);
    } catch (err) {
      console.error("Copy failed", err);
    }
  };
  return (
    <SpeedDial
      ariaLabel="Share post"
      icon={<SpeedDialIcon icon={<ShareIcon />} />}
      direction="down"
      sx={{
        position: "fixed",
        top: "40%",
        right: 24,
      }}
    >
      <SpeedDialAction
        icon={isCopied ? <CheckIcon /> : <FacebookIcon />}
        tooltipTitle={"Copy"}
        onClick={handleCopy}
      />
      {/* {actions.map((action) => (
        <SpeedDialAction
          key={action.name}
          icon={action.icon}
          tooltipTitle={action.name}
          onClick={action.onClick}
        />
      ))} */}
    </SpeedDial>
  );
}
