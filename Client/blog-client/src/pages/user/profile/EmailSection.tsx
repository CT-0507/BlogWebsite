import type { User } from "@/context/AuthContext";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import InputLabel from "@mui/material/InputLabel";
import Switch from "@mui/material/Switch";
import TextField from "@mui/material/TextField";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useRef, useState } from "react";
import { changeEmailSchema, type ChangeEmailFormValues } from "./model/schema";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { changeEmailRequest, sendEmailCode } from "@/api/profileApi";
import CheckIcon from "@mui/icons-material/Check";

const COOLDOWN_SECONDS = 60;

export default function EmailSection() {
  const queryClient = useQueryClient();
  const user = queryClient.getQueryData(["me"]) as User;
  const [editSwitch, setEditSwitch] = useState(false);
  const [cooldown, setCooldown] = useState(0);

  const {
    register,
    handleSubmit,
    resetField,
    setValue,
    getValues,
    formState: { errors, isSubmitting, isDirty, dirtyFields },
  } = useForm<ChangeEmailFormValues>({
    resolver: zodResolver(changeEmailSchema),
    defaultValues: {
      email: user?.email || "",
      confirmCode: "",
    },
    mode: "all",
  });

  const resetForm = () => {
    resetField("email", {
      defaultValue: user?.email || "asd@gm.vf",
    });
    resetField("confirmCode", {
      defaultValue: "",
    });
  };
  const { mutate, isPending } = useMutation({
    mutationFn: changeEmailRequest,
    retry: false,
    onSuccess: (data) => {
      console.log(data);
      // Update me data on login response instead of fetching again
      console.log();
      queryClient.setQueryData(["me"], (oldData: User) => {
        return {
          ...oldData,
          email: data.email,
        };
      });
    },
    onError: (error) => {
      console.error(error.message);
    },
  });
  const intervalRef = useRef<number | null>(null);
  // start cooldown after successful send
  const startCooldown = () => {
    setCooldown(COOLDOWN_SECONDS);

    intervalRef.current = window.setInterval(() => {
      setCooldown((prev) => {
        if (prev <= 1) {
          if (intervalRef.current) {
            clearInterval(intervalRef.current);
            intervalRef.current = null;
          }
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
  };
  const { mutate: sendMutation, isPending: sendIsPending } = useMutation({
    mutationFn: sendEmailCode,
    retry: false,
    onSuccess: (data) => {
      setValue("confirmCode", data.code);
      startCooldown();
    },
    onError: (error) => {
      console.error(error.message);
    },
  });

  const handleCancel = () => {
    resetForm();
  };

  const handleSwitchMode = () => {
    resetForm();
    setEditSwitch((prev) => !prev);
  };
  const handleSendCode = async () => {
    if (errors.email) {
      alert("Email is invalid");
      return;
    }
    sendMutation(getValues("email"));
  };
  const onSubmit = async (data: ChangeEmailFormValues) => {
    console.log("Form Data:", data);
    mutate(data);
  };
  return (
    <>
      <Box
        component="section"
        id="userEmailInfo"
        onSubmit={handleSubmit(onSubmit)}
        sx={{ p: 2 }}
      >
        <Box display="flex" alignItems="center" justifyContent="flex-end">
          <InputLabel htmlFor="editEmailSwitch">Edit</InputLabel>
          <Switch
            id="editEmailSwitch"
            checked={editSwitch}
            onChange={handleSwitchMode}
            disabled={isSubmitting || isPending}
            slotProps={{ input: { "aria-label": "controlled" } }}
          />
        </Box>
        <Box component="form">
          <Box>
            <InputLabel htmlFor="email" sx={{ mb: 1 }}>
              Email address
            </InputLabel>
            <TextField
              id="email"
              disabled={!editSwitch || isSubmitting || isPending}
              placeholder="Enter your email"
              {...register("email")}
              size="small"
              error={!!errors.email}
              helperText={errors.email?.message || " "}
              sx={{ width: "45%", mr: 1 }}
            />
            <CheckIcon
              sx={{
                color: "green",
                visibility: dirtyFields.email ? "visible" : "hidden",
                mr: 1,
              }}
              titleAccess="Modified field"
            />
            <Button
              variant="contained"
              onClick={handleSendCode}
              disabled={
                !editSwitch ||
                isSubmitting ||
                isPending ||
                !!errors.email ||
                sendIsPending ||
                cooldown > 0
              }
            >
              {isPending && "Sending…"}
              {!isPending && cooldown > 0 && `Resend in ${cooldown}s`}
              {!isPending && cooldown === 0 && "Send Email"}
            </Button>
          </Box>
          <InputLabel htmlFor="code" sx={{ mb: 1 }}>
            Code
          </InputLabel>
          <TextField
            id="code"
            placeholder="Enter confirm code"
            {...register("confirmCode")}
            size="small"
            error={!!errors.confirmCode}
            disabled={!editSwitch || isSubmitting || isPending}
            sx={{ width: "45%" }}
            helperText={errors.confirmCode?.message || " "}
          />
          <Box id="email-form-action" sx={{ mt: 1 }}>
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
