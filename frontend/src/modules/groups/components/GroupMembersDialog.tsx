import { useEffect, useMemo, useState } from 'react';
import {
  Alert,
  Autocomplete,
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  List,
  ListItem,
  ListItemText,
  MenuItem,
  Stack,
  TextField,
  Tooltip,
  Typography,
} from '@mui/material';
import DeleteOutlineIcon from '@mui/icons-material/DeleteOutline';
import type { SelectChangeEvent } from '@mui/material/Select';
import { LoadingSpinner } from '../../../components';
import { useApi } from '../../../hooks';
import { API_ENDPOINTS, GROUP_ROLES } from '../../../constants';
import { useLanguage } from '../../../contexts/LanguageContext';
import type {
  ApiError,
  GroupItem,
  GroupMember,
  GroupRole,
  UsersResponse,
} from '../../../types';

interface GroupMembersDialogProps {
  open: boolean;
  group: GroupItem | null;
  members: GroupMember[];
  loading: boolean;
  mutating: boolean;
  error: ApiError | null;
  onClose: () => void;
  onAddMember: (userUUID: string, role: GroupRole) => Promise<void>;
  onUpdateMemberRole: (userUUID: string, role: GroupRole) => Promise<void>;
  onRemoveMember: (userUUID: string) => Promise<void>;
}

const defaultRole: GroupRole = GROUP_ROLES.REVIEWER;

const GroupMembersDialog = ({
  open,
  group,
  members,
  loading,
  mutating,
  error,
  onClose,
  onAddMember,
  onUpdateMemberRole,
  onRemoveMember,
}: GroupMembersDialogProps) => {
  const { t } = useLanguage();
  const [selectedUser, setSelectedUser] = useState<string>('');
  const [selectedRole, setSelectedRole] = useState<GroupRole>(defaultRole);
  const [localError, setLocalError] = useState<string | null>(null);

  const {
    data: usersData,
    loading: usersLoading,
    error: usersError,
    refetch: refetchUsers,
  } = useApi<UsersResponse>(API_ENDPOINTS.USERS.BASE, {
    immediate: open,
  });

  useEffect(() => {
    if (!open) {
      setSelectedUser('');
      setSelectedRole(defaultRole);
      setLocalError(null);
    }
  }, [open]);

  useEffect(() => {
    if (open && !usersData && !usersLoading) {
      void refetchUsers();
    }
  }, [open, usersData, usersLoading, refetchUsers]);

  const availableUsers = useMemo(() => {
    if (!usersData?.users) {
      return [] as Array<{ uuid: string; login: string; email: string }>;
    }

    const memberUUIDs = new Set(members.map((member) => member.user_uuid));
    return usersData.users.filter((user) => !memberUUIDs.has(user.uuid));
  }, [usersData?.users, members]);

  const handleAddMember = async () => {
    if (!selectedUser) {
      setLocalError(t('groups.membersAddValidation'));
      return;
    }

    try {
      await onAddMember(selectedUser, selectedRole);
      setSelectedUser('');
      setSelectedRole(defaultRole);
      setLocalError(null);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : t('groups.membersActionError');
      setLocalError(message);
    }
  };

  const handleRoleChange =
    (userUUID: string) => async (event: SelectChangeEvent<GroupRole>) => {
      const newRole = event.target.value as GroupRole;
      try {
        await onUpdateMemberRole(userUUID, newRole);
        setLocalError(null);
      } catch (err) {
        const message =
          err instanceof Error ? err.message : t('groups.membersActionError');
        setLocalError(message);
      }
    };

  const handleRemoveMember = (userUUID: string) => async () => {
    try {
      await onRemoveMember(userUUID);
      setLocalError(null);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : t('groups.membersActionError');
      setLocalError(message);
    }
  };

  const roleOptions = useMemo(
    () => [
      { value: GROUP_ROLES.COAUTHOR, label: t('groups.role.coauthor') },
      { value: GROUP_ROLES.REVIEWER, label: t('groups.role.reviewer') },
    ],
    [t],
  );

  const resolveRoleName = (role: GroupRole) => {
    if (role === GROUP_ROLES.AUTHOR) {
      return t('groups.role.author');
    }
    const match = roleOptions.find((option) => option.value === role);
    return match?.label ?? role;
  };

  const dialogTitle = group
    ? t('groups.membersTitle').replace('{name}', group.name)
    : t('groups.membersTitleFallback');

  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth='sm'>
      <DialogTitle>{dialogTitle}</DialogTitle>
      <DialogContent dividers>
        <Stack spacing={3}>
          <Typography variant='body2' color='text.secondary'>
            {t('groups.membersSubtitle')}
          </Typography>

          {(error || localError) && (
            <Alert severity='error'>{localError ?? error?.message}</Alert>
          )}

          {usersError && !usersLoading && (
            <Alert severity='warning'>
              {usersError.message || t('groups.membersUsersError')}
            </Alert>
          )}

          <Stack spacing={2} direction={{ xs: 'column', md: 'row' }}>
            <Autocomplete
              fullWidth
              size='small'
              disabled={mutating || usersLoading}
              options={availableUsers}
              value={
                availableUsers.find((option) => option.uuid === selectedUser) ??
                null
              }
              onChange={(_, option) => setSelectedUser(option?.uuid ?? '')}
              getOptionLabel={(option) => `${option.login} (${option.email})`}
              loading={usersLoading}
              renderInput={(params) => (
                <TextField
                  {...params}
                  label={t('groups.membersAddUser')}
                  placeholder={t('groups.membersAddUserPlaceholder')}
                />
              )}
            />

            <TextField
              select
              size='small'
              label={t('groups.membersRoleLabel')}
              value={selectedRole}
              onChange={(event) =>
                setSelectedRole(event.target.value as GroupRole)
              }
              sx={{ minWidth: 160 }}
              disabled={mutating}
            >
              {roleOptions.map((option) => (
                <MenuItem key={option.value} value={option.value}>
                  {option.label}
                </MenuItem>
              ))}
            </TextField>

            <Button
              variant='contained'
              onClick={handleAddMember}
              disabled={!selectedUser || mutating}
            >
              {t('groups.membersAddButton')}
            </Button>
          </Stack>

          {loading ? (
            <LoadingSpinner py={4} />
          ) : members.length === 0 ? (
            <Typography variant='body2' color='text.secondary'>
              {t('groups.membersEmpty')}
            </Typography>
          ) : (
            <List disablePadding>
              {members.map((member) => {
                const isAuthor = member.role === GROUP_ROLES.AUTHOR;
                const roleLabel = resolveRoleName(member.role);

                return (
                  <ListItem
                    key={member.user_uuid}
                    sx={{
                      border: '1px solid',
                      borderColor: 'divider',
                      borderRadius: 2,
                      mb: 1.5,
                    }}
                    secondaryAction={
                      isAuthor ? (
                        <Typography variant='body2' color='text.secondary'>
                          {roleLabel}
                        </Typography>
                      ) : (
                        <Stack direction='row' spacing={1} alignItems='center'>
                          <TextField
                            select
                            size='small'
                            value={member.role}
                            onChange={handleRoleChange(member.user_uuid)}
                            disabled={mutating}
                          >
                            {roleOptions.map((option) => (
                              <MenuItem key={option.value} value={option.value}>
                                {option.label}
                              </MenuItem>
                            ))}
                          </TextField>
                          <Tooltip title={t('groups.membersRemoveTooltip')}>
                            <span>
                              <IconButton
                                size='small'
                                color='error'
                                onClick={handleRemoveMember(member.user_uuid)}
                                disabled={mutating}
                              >
                                <DeleteOutlineIcon fontSize='small' />
                              </IconButton>
                            </span>
                          </Tooltip>
                        </Stack>
                      )
                    }
                  >
                    <ListItemText
                      primary={member.user_login}
                      secondary={
                        <Typography variant='caption' color='text.secondary'>
                          {member.user_email}
                        </Typography>
                      }
                    />
                  </ListItem>
                );
              })}
            </List>
          )}
        </Stack>
      </DialogContent>
      <DialogActions>
        <Box sx={{ flex: 1 }} />
        <Button onClick={onClose}>{t('groups.cancel')}</Button>
      </DialogActions>
    </Dialog>
  );
};

export default GroupMembersDialog;
