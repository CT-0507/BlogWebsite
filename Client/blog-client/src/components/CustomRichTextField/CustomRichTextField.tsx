import TextField from "@mui/material/TextField";
import { Controller, useForm } from "react-hook-form";
import { FormControl } from "@mui/material";

type FormValues = {
  message: string;
};

export default function FormExample() {
  const maxLength = 280;

  const { control, handleSubmit } = useForm<FormValues>({
    defaultValues: {
      message: "",
    },
  });

  const onSubmit = (data: FormValues) => {
    console.log(data);
  };

  return (
    <FormControl
      component="form"
      onSubmit={handleSubmit(onSubmit)}
      sx={{ width: 500 }}
    >
      <Controller
        name="message"
        control={control}
        rules={{
          required: "Message is required",
          maxLength: {
            value: maxLength,
            message: `Maximum ${maxLength} characters`,
          },
        }}
        render={({ field, fieldState }) => (
          <TextField
            {...field}
            fullWidth
            multiline
            minRows={4}
            label="Message"
            variant="outlined"
            error={!!fieldState.error}
            helperText={
              fieldState.error?.message ||
              `${field.value?.length || 0}/${maxLength}`
            }
          />
        )}
      />
    </FormControl>
  );
}
