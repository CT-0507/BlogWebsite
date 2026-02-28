import { useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { signupSchema, type SignupFormValues } from "./model/schema";
import { zodResolver } from "@hookform/resolvers/zod";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Box from "@mui/material/Box";
import FormControl from "@mui/material/FormControl";
import InputLabel from "@mui/material/InputLabel";
import OutlinedInput from "@mui/material/OutlinedInput";
import InputAdornment from "@mui/material/InputAdornment";
import IconButton from "@mui/material/IconButton";
import VisibilityOff from "@mui/icons-material/VisibilityOff";
import Visibility from "@mui/icons-material/Visibility";
import FormHelperText from "@mui/material/FormHelperText";
import { useMutation } from "@tanstack/react-query";
import { signupRequest } from "@/api/auth";
import { useNavigate } from "react-router-dom";

export default function SignupForm() {
  const navigate = useNavigate();
  const {
    register,
    control,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<SignupFormValues>({
    resolver: zodResolver(signupSchema),
    defaultValues: {
      username: "",
      password: "",
      confirmPassword: "",
    },
    mode: "all",
  });

  const [showPassword, setShowPassword] = useState(false);

  const handleClickShowPassword = () => setShowPassword((show) => !show);

  const handleMouseDownPassword = (
    event: React.MouseEvent<HTMLButtonElement>
  ) => {
    event.preventDefault();
    setShowPassword(true);
  };

  const handleMouseUpPassword = (
    event: React.MouseEvent<HTMLButtonElement>
  ) => {
    event.preventDefault();
    setShowPassword(false);
  };

  const { mutate, isPending } = useMutation({
    mutationFn: signupRequest,
    onSuccess: (data) => {
      console.log("Logged in:", data);
      navigate(`/home`);
      // save token / redirect
    },
    onError: (error) => {
      // console.error(error.response?.data?.message || error.message);
      console.error(error.message);
    },
  });

  const onSubmit = async (data: SignupFormValues) => {
    console.log("Form Data:", data);
    mutate(data);
  };
  return (
    <>
      <Typography variant="h5" mb={3} textAlign="center">
        Sign Up for new account now
      </Typography>

      <Box
        flex={1}
        sx={{
          display: "flex",
          flexDirection: "column",
          justifyContent: "space-between",
        }}
        component="form"
        onSubmit={handleSubmit(onSubmit)}
      >
        <Box id="form-field">
          <TextField
            label="Username"
            fullWidth
            margin="normal"
            {...register("username")}
            error={!!errors.username}
            helperText={errors.username?.message || " "}
          />

          <FormControl fullWidth>
            <InputLabel>Password</InputLabel>
            <Controller
              control={control}
              name="password"
              render={({ field }) => (
                <OutlinedInput
                  id="outlined-adornment-password"
                  type={showPassword ? "text" : "password"}
                  {...field}
                  error={!!errors.password}
                  fullWidth
                  endAdornment={
                    <InputAdornment position="end">
                      <IconButton
                        aria-label={
                          showPassword
                            ? "hide the password"
                            : "display the password"
                        }
                        onClick={handleClickShowPassword}
                        onMouseDown={handleMouseDownPassword}
                        onMouseUp={handleMouseUpPassword}
                        edge="end"
                      >
                        {showPassword ? <VisibilityOff /> : <Visibility />}
                      </IconButton>
                    </InputAdornment>
                  }
                  label="Password"
                />
              )}
            />
            <FormHelperText error={!!errors.password}>
              {errors.password?.message || " "}
            </FormHelperText>
          </FormControl>

          <TextField
            label="Confirm Password"
            fullWidth
            type="password"
            margin="normal"
            {...register("confirmPassword")}
            error={!!errors.confirmPassword}
            helperText={errors.confirmPassword?.message || " "}
          />

          <Box display={"flex"} justifyContent={"space-between"}>
            <TextField
              label="First name"
              margin="normal"
              {...register("firstName")}
              error={!!errors.firstName}
              helperText={errors.firstName?.message || " "}
            />

            <TextField
              label="Last name Name"
              margin="normal"
              {...register("lastName")}
              error={!!errors.lastName}
              helperText={errors.lastName?.message || " "}
            />
          </Box>
        </Box>

        <Button
          type="submit"
          variant="contained"
          fullWidth
          sx={{ mt: 3 }}
          disabled={isSubmitting || isPending}
        >
          {isPending ? "Signing up..." : "Sign Up"}
        </Button>
      </Box>
    </>
  );
}
