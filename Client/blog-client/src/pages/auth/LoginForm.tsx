import { loginSchema, type LoginFormValues } from "./model/schema";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { useState } from "react";
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
import { useNavigate } from "react-router-dom";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { loginRequest } from "@/api/auth";
import { tokenStore } from "@/api/store/tokenStore";
import { useAuth } from "@/hooks/useAuth";

export default function LoginForm() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const {
    register,
    control,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      username: "",
      password: "",
    },
    mode: "all",
  });
  const auth = useAuth();
  console.log(auth);

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
    mutationFn: loginRequest,
    retry: false,
    onSuccess: (data) => {
      console.log("Logged in:", data);
      tokenStore.set(data.accessToken);
      queryClient.setQueryData(["me"], data);
      navigate(`/dashboard`);
      // save token / redirect
    },
    onError: (error) => {
      // console.error(error.response?.data?.message || error.message);
      console.error(error.message);
    },
  });

  const onSubmit = async (data: LoginFormValues) => {
    console.log("Form Data:", data);
    // Call API here
    mutate(data);
  };
  return (
    <>
      <Typography variant="h5" mb={3} textAlign="center">
        This is login form
      </Typography>

      <Box
        flex={1}
        component="form"
        sx={{
          display: "flex",
          flexDirection: "column",
          justifyContent: "space-between",
        }}
        onSubmit={handleSubmit(onSubmit)}
      >
        <Box
          sx={{
            p: 0,
          }}
          id="form-field"
        >
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
        </Box>

        <Button
          type="submit"
          variant="contained"
          fullWidth
          sx={{ mt: 3 }}
          disabled={isSubmitting || isPending}
        >
          {isPending ? "Logging in..." : "Login"}
        </Button>
      </Box>
    </>
  );
}
