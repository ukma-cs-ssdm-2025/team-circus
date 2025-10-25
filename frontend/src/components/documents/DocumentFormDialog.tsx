import {
  useEffect,
  useMemo,
  useState,
  type ChangeEvent,
  type FormEvent,
} from 'react';
import {
  Alert,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  MenuItem,
  Stack,
  TextField,
} from '@mui/material';
import type { BaseComponentProps } from '../../types';
import type { GroupOption } from '../../types/entities';

interface DocumentFormDialogProps extends BaseComponentProps {
  open: boolean;
  title: string;
  confirmLabel: string;
  cancelLabel: string;
  nameLabel: string;
  namePlaceholder: string;
  nameHelperText: string;
  groupLabel: string;
  groupOptions: GroupOption[];
  loading?: boolean;
  errorMessage?: string | null;
  onClose: () => void;
  onSubmit: (payload: {
    name: string;
    groupUUID: string;
  }) => void | Promise<void>;
}

const DocumentFormDialog = ({
  open,
  title,
  confirmLabel,
  cancelLabel,
  nameLabel,
  namePlaceholder,
  nameHelperText,
  groupLabel,
  groupOptions,
  loading = false,
  errorMessage = null,
  onClose,
  onSubmit,
  className = '',
}: DocumentFormDialogProps) => {
  const [name, setName] = useState('');
  const [selectedGroup, setSelectedGroup] = useState('');
  const [touched, setTouched] = useState(false);

  const defaultGroup = useMemo(
    () => groupOptions[0]?.value ?? '',
    [groupOptions],
  );

  useEffect(() => {
    if (open) {
      setName('');
      setSelectedGroup(defaultGroup);
      setTouched(false);
    }
  }, [defaultGroup, open]);

  const handleNameChange = (event: ChangeEvent<HTMLInputElement>) => {
    if (!touched) {
      setTouched(true);
    }
    setName(event.target.value);
  };

  const handleGroupChange = (event: ChangeEvent<HTMLInputElement>) => {
    setSelectedGroup(event.target.value);
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setTouched(true);

    const trimmedName = name.trim();

    if (!trimmedName || !selectedGroup || loading) {
      return;
    }

    await onSubmit({
      name: trimmedName,
      groupUUID: selectedGroup,
    });
  };

  const isNameValid = name.trim().length > 0;
  const isFormValid = isNameValid && Boolean(selectedGroup);

  return (
    <Dialog
      open={open}
      onClose={loading ? undefined : onClose}
      disableEscapeKeyDown={loading}
      fullWidth
      maxWidth='sm'
      className={className}
      component='form'
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
            onChange={handleNameChange}
            disabled={loading}
            error={touched && !isNameValid}
            helperText={touched && !isNameValid ? nameHelperText : ' '}
          />

          <TextField
            select
            fullWidth
            label={groupLabel}
            value={selectedGroup}
            onChange={handleGroupChange}
            disabled={loading || groupOptions.length === 0}
          >
            {groupOptions.map((option) => (
              <MenuItem key={option.value} value={option.value}>
                {option.label}
              </MenuItem>
            ))}
          </TextField>

          {errorMessage && <Alert severity='error'>{errorMessage}</Alert>}
        </Stack>
      </DialogContent>
      <DialogActions sx={{ px: 3, pb: 3 }}>
        <Button onClick={onClose} disabled={loading}>
          {cancelLabel}
        </Button>
        <Button
          type='submit'
          variant='contained'
          disabled={!isFormValid || loading}
        >
          {loading ? `${confirmLabel}...` : confirmLabel}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default DocumentFormDialog;
