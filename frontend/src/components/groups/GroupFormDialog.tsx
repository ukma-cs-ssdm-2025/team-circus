import { useEffect, useState, type ChangeEvent, type FormEvent } from "react";
import {
  Alert,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Stack,
  TextField,
} from "@mui/material";
import type { BaseComponentProps } from "../../types";

interface GroupFormDialogProps extends BaseComponentProps {
  open: boolean;
  title: string;
  confirmLabel: string;
  cancelLabel: string;
  nameLabel: string;
  namePlaceholder: string;
  nameHelperText: string;
  initialName?: string;
  loading?: boolean;
  errorMessage?: string | null;
  onClose: () => void;
  onSubmit: (name: string) => void | Promise<void>;
}

const GroupFormDialog = ({
  open,
  title,
  confirmLabel,
  cancelLabel,
  nameLabel,
  namePlaceholder,
  nameHelperText,
  initialName = "",
  loading = false,
  errorMessage = null,
  onClose,
  onSubmit,
  className = "",
}: GroupFormDialogProps) => {
  const [name, setName] = useState(initialName);
  const [touched, setTouched] = useState(false);

  useEffect(() => {
    if (open) {
      setName(initialName);
      setTouched(false);
    }
  }, [initialName, open]);

  const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
    if (!touched) {
      setTouched(true);
    }
    setName(event.target.value);
  };

  const handleSubmit = async (event: FormEvent<HTMLDivElement>) => {
    event.preventDefault();
    setTouched(true);

    if (!name.trim() || loading) {
      return;
    }

    await onSubmit(name.trim());
  };

  const isNameValid = name.trim().length > 0;

  return (
    <Dialog
      open={open}
      onClose={loading ? undefined : onClose}
      disableEscapeKeyDown={loading}
      fullWidth
      maxWidth="sm"
      className={className}
      component="form"
      onSubmit={handleSubmit}
    >
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        <Stack spacing={2} sx={{ mt: 1 }}>
          <TextField
            autoFocus
            fullWidth
            label={nameLabel}
            placeholder={namePlaceholder}
            value={name}
            onChange={handleChange}
            disabled={loading}
            error={touched && !isNameValid}
            helperText={touched && !isNameValid ? nameHelperText : " "}
          />

          {errorMessage && <Alert severity="error">{errorMessage}</Alert>}
        </Stack>
      </DialogContent>
      <DialogActions sx={{ px: 3, pb: 3 }}>
        <Button onClick={onClose} disabled={loading}>
          {cancelLabel}
        </Button>
        <Button
          type="submit"
          variant="contained"
          disabled={!isNameValid || loading}
        >
          {loading ? `${confirmLabel}...` : confirmLabel}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default GroupFormDialog;
