import type { User } from "@/context/AuthContext";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import InputLabel from "@mui/material/InputLabel";
import Switch from "@mui/material/Switch";
import TextField from "@mui/material/TextField";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import {
  changeBasicInfoSchema,
  type ChangeBasicInfoFormValues,
} from "./model/schema";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { changeBasicInfoRequest } from "@/api/profileApi";
import CheckIcon from "@mui/icons-material/Check";

export default function BasicInfoSection() {
  const queryClient = useQueryClient();
  const user = queryClient.getQueryData(["me"]) as User;
  const [editSwitch, setEditSwitch] = useState(false);

  const {
    register,
    handleSubmit,
    resetField,
    formState: { errors, isSubmitting, isDirty, dirtyFields },
  } = useForm<ChangeBasicInfoFormValues>({
    resolver: zodResolver(changeBasicInfoSchema),
    defaultValues: {
      firstName: user.firstName || "",
      lastName: user.lastName || "",
    },
    mode: "all",
  });

  const resetForm = () => {
    resetField("firstName", {
      defaultValue: user?.firstName || "",
    });
    resetField("lastName", {
      defaultValue: user?.lastName || "",
    });
  };
  const { mutate, isPending } = useMutation({
    mutationFn: changeBasicInfoRequest,
    retry: false,
    onSuccess: (data) => {
      // Update me data on login response instead of fetching again
      queryClient.setQueryData(["me"], (oldData: User) => ({
        ...oldData,
        firstName: data?.firstName,
        lastName: data?.lastName,
      }));
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

  const onSubmit = async (data: ChangeBasicInfoFormValues) => {
    console.log("Form Data:", data);
    mutate(data);
  };
  return (
    <>
      <Box
        component="section"
        id="userBasicInfo"
        onSubmit={handleSubmit(onSubmit)}
        sx={{ p: 2 }}
      >
        <Box display="flex" alignItems="center" justifyContent="flex-end">
          <InputLabel htmlFor="editBasicInfoSwitch">Edit</InputLabel>
          <Switch
            id="editBasicInfoSwitch"
            checked={editSwitch}
            onChange={handleSwitchMode}
            disabled={isSubmitting || isPending}
            slotProps={{ input: { "aria-label": "controlled" } }}
          />
        </Box>
        <Box component="form">
          <Box display="flex" justifyContent="space-between">
            <Box sx={{ width: "45%" }}>
              <InputLabel htmlFor="firstName" sx={{ mb: 1 }}>
                First Name
              </InputLabel>
              <TextField
                id="firstName"
                disabled={!editSwitch || isSubmitting || isPending}
                placeholder="Enter your first name"
                {...register("firstName")}
                size="small"
                error={!!errors.firstName}
                helperText={errors.firstName?.message || " "}
                sx={{ mr: 1, width: "90%" }}
              />
              <CheckIcon
                sx={{
                  display: dirtyFields.firstName ? "inline" : "none",
                  color: "green",
                }}
                titleAccess="Modified field"
              />
            </Box>
            <Box sx={{ width: "40%" }}>
              <InputLabel htmlFor="lastName" sx={{ mb: 1 }}>
                Last Name
              </InputLabel>
              <TextField
                id="lastName"
                placeholder="Enter your last name"
                {...register("lastName")}
                size="small"
                error={!!errors.lastName}
                disabled={!editSwitch || isSubmitting || isPending}
                helperText={errors.lastName?.message || " "}
                sx={{ width: "90%" }}
              />
              <CheckIcon
                sx={{
                  display: dirtyFields.lastName ? "inline" : "none",
                  color: "green",
                }}
                titleAccess="Modified field"
              />
            </Box>
          </Box>
          <Box id="basic-info-form-action" sx={{ mt: 1 }}>
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
