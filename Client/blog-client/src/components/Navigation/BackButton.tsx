import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { useNavigate } from "react-router-dom";

export default function BackButton({ fallbackPath = "/my-blogs" }) {
  const navigate = useNavigate();

  const handleBack = () => {
    if (window.history.length > 1) {
      navigate(-1);
    } else {
      navigate(fallbackPath);
    }
  };

  return (
    <Tooltip title="Back">
      <IconButton onClick={handleBack}>
        <ArrowBackIcon />
      </IconButton>
    </Tooltip>
  );
}
