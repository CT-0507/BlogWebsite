import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import InputLabel from "@mui/material/InputLabel";
import Switch from "@mui/material/Switch";
import TextField from "@mui/material/TextField";
import { useMutation } from "@tanstack/react-query";
import { useState } from "react";
import {
  changePasswordSchema,
  type ChangePasswordFormValues,
} from "./model/schema";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { changePasswordRequest } from "@/api/profileApi";
import FormControl from "@mui/material/FormControl";
import OutlinedInput from "@mui/material/OutlinedInput";
import InputAdornment from "@mui/material/InputAdornment";
import IconButton from "@mui/material/IconButton";
import VisibilityOff from "@mui/icons-material/VisibilityOff";
import Visibility from "@mui/icons-material/Visibility";
import FormHelperText from "@mui/material/FormHelperText";

export default function PasswordSection() {
  const [editSwitch, setEditSwitch] = useState(false);
  const [showCurrentPassword, setShowCurrentPassword] = useState(false);
  const [showNewPassword, setShowNewPassword] = useState(false);
  const {
    control,
    register,
    handleSubmit,
    resetField,
    formState: { errors, isSubmitting, isDirty },
  } = useForm<ChangePasswordFormValues>({
    resolver: zodResolver(changePasswordSchema),
    defaultValues: {
      currentPassword: "",
      newPassword: "",
      confirmNewPassword: "",
    },
    mode: "all",
  });

  const resetForm = () => {
    resetField("currentPassword", {
      defaultValue: "",
    });
    resetField("newPassword", {
      defaultValue: "",
    });
    resetField("confirmNewPassword", {
      defaultValue: "",
    });
  };
  const { mutate, isPending } = useMutation({
    mutationFn: changePasswordRequest,
    retry: false,
    onSuccess: (data) => {
      console.log(data);
    },
    onError: (error) => {
      console.error(error.message);
    },
  });

  const handleCancel = () => {
    resetForm();
  };

  const handleMouseDownCurrentPassword = (
    event: React.MouseEvent<HTMLButtonElement>
  ) => {
    event.preventDefault();
    setShowCurrentPassword(true);
  };

  const handleMouseUpCurrentPassword = (
    event: React.MouseEvent<HTMLButtonElement>
  ) => {
    event.preventDefault();
    setShowCurrentPassword(false);
  };

  const handleMouseDownNewPassword = (
    event: React.MouseEvent<HTMLButtonElement>
  ) => {
    event.preventDefault();
    setShowNewPassword(true);
  };

  const handleMouseUpNewPassword = (
    event: React.MouseEvent<HTMLButtonElement>
  ) => {
    event.preventDefault();
    setShowNewPassword(false);
  };

  const handleSwitchMode = () => {
    resetForm();
    setEditSwitch((prev) => !prev);
  };
  const onSubmit = async (data: ChangePasswordFormValues) => {
    console.log("Form Data:", data);
    mutate(data);
  };
  return (
    <>
      <Box
        component="section"
        id="userPassword"
        onSubmit={handleSubmit(onSubmit)}
        sx={{ p: 2 }}
      >
        <Box display="flex" alignItems="center" justifyContent="flex-end">
          <InputLabel htmlFor="editPasswordSwitch">Edit</InputLabel>
          <Switch
            id="editPasswordSwitch"
            checked={editSwitch}
            onChange={handleSwitchMode}
            disabled={isSubmitting || isPending}
            slotProps={{ input: { "aria-label": "controlled" } }}
          />
        </Box>
        <Box component="form">
          <Box sx={{ width: "45%", mr: 1 }}>
            <InputLabel htmlFor="password">Current Password</InputLabel>
            <FormControl fullWidth>
              <Controller
                control={control}
                name="currentPassword"
                disabled={!editSwitch}
                render={({ field }) => (
                  <OutlinedInput
                    id="password"
                    type={showCurrentPassword ? "text" : "password"}
                    {...field}
                    error={!!errors.currentPassword}
                    fullWidth
                    size="small"
                    endAdornment={
                      <InputAdornment position="end">
                        <IconButton
                          disabled={!editSwitch}
                          aria-label={
                            showCurrentPassword
                              ? "hide the password"
                              : "display the password"
                          }
                          onMouseDown={handleMouseDownCurrentPassword}
                          onMouseUp={handleMouseUpCurrentPassword}
                          edge="end"
                        >
                          {showCurrentPassword ? (
                            <VisibilityOff />
                          ) : (
                            <Visibility />
                          )}
                        </IconButton>
                      </InputAdornment>
                    }
                  />
                )}
              />
              <FormHelperText error={!!errors.currentPassword}>
                {errors.currentPassword?.message || " "}
              </FormHelperText>
            </FormControl>
          </Box>
          <Box sx={{ width: "45%", mr: 1 }}>
            <InputLabel htmlFor="newPassword">New Password</InputLabel>
            <FormControl fullWidth>
              <Controller
                control={control}
                name="newPassword"
                disabled={!editSwitch}
                render={({ field }) => (
                  <OutlinedInput
                    id="newPassword"
                    type={showNewPassword ? "text" : "password"}
                    {...field}
                    error={!!errors.newPassword}
                    fullWidth
                    size="small"
                    endAdornment={
                      <InputAdornment position="end">
                        <IconButton
                          disabled={!editSwitch}
                          aria-label={
                            showNewPassword
                              ? "hide the new password"
                              : "display the new password"
                          }
                          onMouseDown={handleMouseDownNewPassword}
                          onMouseUp={handleMouseUpNewPassword}
                          edge="end"
                        >
                          {showNewPassword ? <VisibilityOff /> : <Visibility />}
                        </IconButton>
                      </InputAdornment>
                    }
                  />
                )}
              />
              <FormHelperText error={!!errors.newPassword}>
                {errors.newPassword?.message || " "}
              </FormHelperText>
            </FormControl>
          </Box>
          <Box>
            <InputLabel htmlFor="confirmNewPassword" sx={{ mb: 1 }}>
              Confirm New Password
            </InputLabel>
            <TextField
              id="confirmNewPassword"
              placeholder="Enter Confirm New Password"
              {...register("confirmNewPassword")}
              type="password"
              size="small"
              error={!!errors.confirmNewPassword}
              disabled={!editSwitch || isSubmitting || isPending}
              sx={{ width: "45%" }}
              helperText={errors.confirmNewPassword?.message || " "}
            />
          </Box>
          <Box id="password-form-action" sx={{ mt: 1 }}>
            <Button
              variant="contained"
              type="submit"
              disabled={!editSwitch || !isDirty}
              sx={{
                width: "150px",
                mr: 1,
              }}
            >
              {isPending ? "Saving..." : "Save"}
            </Button>
            <Button
              variant="contained"
              onClick={handleCancel}
              disabled={!editSwitch || isSubmitting || isPending}
              color="error"
              sx={{
                width: "150px",
              }}
            >
              Cancel
            </Button>
          </Box>
        </Box>
      </Box>
    </>
  );
}
