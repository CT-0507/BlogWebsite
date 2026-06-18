import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Stack from "@mui/material/Stack";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import IconButton from "@mui/material/IconButton";
import Snackbar from "@mui/material/Snackbar";
import { contactSchema, type ContactFormValues } from "./model/schema";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { useCreateContact } from "@/hooks/useContact";
import CloseIcon from "@mui/icons-material/Close";
import { useState } from "react";

interface ContactFormProps {
  id?: string;
}

export default function ContactForm({ id }: ContactFormProps) {
  const [open, setOpen] = useState(false);
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<ContactFormValues>({
    resolver: zodResolver(contactSchema),
    defaultValues: {
      email: "",
      content: "",
    },
    mode: "all",
  });

  const { mutate, isPending } = useCreateContact();

  const onSubmit = async (data: ContactFormValues) => {
    mutate(data, {
      onSuccess: () => {
        reset();
        setOpen(true);
      },
    });
  };

  const handleCloseSnackbar = () => {
    setOpen(false);
  };

  return (
    <Box
      minHeight="40vh"
      display="flex"
      alignItems="center"
      id={id}
      component="form"
      sx={{ mb: 2 }}
      onSubmit={handleSubmit(onSubmit)}
    >
      <Snackbar
        open={open}
        autoHideDuration={6000}
        onClose={handleCloseSnackbar}
        message="Successfully created a contact request"
        action={
          <>
            <IconButton
              size="small"
              aria-label="close"
              color="inherit"
              onClick={handleCloseSnackbar}
            >
              <CloseIcon fontSize="small" />
            </IconButton>
          </>
        }
      />
      <Stack spacing={2} sx={{ width: "100%" }}>
        <Typography variant="h3">Let's Connect</Typography>

        <Typography>
          Open to collaboration, architecture discussions, and engineering
          opportunities.
        </Typography>

        <TextField
          type="email"
          placeholder="Your email address"
          {...register("email")}
          error={!!errors.email}
          helperText={errors.email?.message || " "}
        />

        <TextField
          multiline
          rows={15}
          sx={{ width: "100%" }}
          placeholder="Write your message here"
          {...register("content")}
          error={!!errors.content}
          helperText={errors.content?.message || " "}
        />

        <Button
          variant="contained"
          type="submit"
          disabled={isSubmitting || isPending}
        >
          Contact
        </Button>
      </Stack>
    </Box>
  );
}
