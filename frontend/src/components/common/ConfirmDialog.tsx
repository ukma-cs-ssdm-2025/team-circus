import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  Stack,
  Alert,
} from '@mui/material';
import type { BaseComponentProps } from '../../types';

interface ConfirmDialogProps extends BaseComponentProps {
  open: boolean;
  title: string;
  description: string;
  confirmLabel: string;
  cancelLabel: string;
  loading?: boolean;
  errorMessage?: string | null;
  onCancel: () => void;
  onConfirm: () => void;
}

const ConfirmDialog = ({
  open,
  title,
  description,
  confirmLabel,
  cancelLabel,
  loading = false,
  errorMessage = null,
  onCancel,
  onConfirm,
  className = '',
}: ConfirmDialogProps) => {
  return (
    <Dialog
      open={open}
      onClose={loading ? undefined : onCancel}
      className={className}
      fullWidth
      maxWidth='xs'
      disableEscapeKeyDown={loading}
    >
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        <Stack spacing={2}>
          <DialogContentText>{description}</DialogContentText>
          {errorMessage && <Alert severity='error'>{errorMessage}</Alert>}
        </Stack>
      </DialogContent>
      <DialogActions sx={{ px: 3, pb: 3 }}>
        <Button onClick={onCancel} disabled={loading}>
          {cancelLabel}
        </Button>
        <Button
          variant='contained'
          color='error'
          onClick={onConfirm}
          disabled={loading}
        >
          {loading ? `${confirmLabel}...` : confirmLabel}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default ConfirmDialog;
