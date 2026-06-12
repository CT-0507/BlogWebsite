import Stepper from "@mui/material/Stepper";
import Step from "@mui/material/Step";
import StepLabel from "@mui/material/StepLabel";

interface Item {
  label: string;
}

interface Props {
  items: Item[];
}

export default function JourneyStepper({ items }: Props) {
  return (
    <Stepper orientation="vertical" activeStep={items.length - 1}>
      {items.map((item) => (
        <Step key={item.label}>
          <StepLabel>{item.label}</StepLabel>
        </Step>
      ))}
    </Stepper>
  );
}
