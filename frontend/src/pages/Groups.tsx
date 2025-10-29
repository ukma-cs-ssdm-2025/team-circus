import { useEffect, useMemo, useState } from 'react';
import { Alert, Button, Snackbar, Stack, Typography } from '@mui/material';
import type { AlertColor } from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import { useNavigate } from 'react-router-dom';
import {
  CenteredContent,
  PageCard,
  PageHeader,
  ErrorAlert,
  LoadingSpinner,
  GroupsList,
  GroupFormDialog,
  ConfirmDialog,
} from '../components';
import { useLanguage } from '../contexts/LanguageContext';
import { useApi, useMutation } from '../hooks';
import { API_ENDPOINTS, GROUP_ROLES, ROUTES } from '../constants';
import { GroupMembersDialog, useGroupMembers } from '../modules/groups';
import type {
  BaseComponentProps,
  GroupsResponse,
  GroupItem as GroupItemType,
  CreateGroupPayload,
  UpdateGroupPayload,
  GroupRole,
} from '../types';

type GroupsProps = BaseComponentProps;

const Groups = ({ className = '' }: GroupsProps) => {
  const { t } = useLanguage();
  const navigate = useNavigate();
  const [formState, setFormState] = useState<{
    mode: 'create' | 'edit';
    group?: GroupItemType;
  } | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<GroupItemType | null>(null);
  const [membersTarget, setMembersTarget] = useState<GroupItemType | null>(
    null,
  );
  const [snackbar, setSnackbar] = useState<{
    message: string;
    severity: AlertColor;
  } | null>(null);

  const { data, loading, error, refetch } = useApi<GroupsResponse>(
    API_ENDPOINTS.GROUPS.BASE,
  );
  const groups = data?.groups ?? [];

  const {
    mutate: createGroup,
    loading: creating,
    error: createError,
    reset: resetCreate,
  } = useMutation<GroupItemType, CreateGroupPayload>(
    API_ENDPOINTS.GROUPS.BASE,
    'POST',
  );

  const updateEndpoint = useMemo(() => {
    if (formState?.mode === 'edit' && formState.group) {
      return `${API_ENDPOINTS.GROUPS.BASE}/${formState.group.uuid}`;
    }
    return '';
  }, [formState]);

  const {
    mutate: updateGroup,
    loading: updating,
    error: updateError,
    reset: resetUpdate,
  } = useMutation<GroupItemType, UpdateGroupPayload>(updateEndpoint, 'PUT');

  const deleteEndpoint = useMemo(() => {
    if (deleteTarget) {
      return `${API_ENDPOINTS.GROUPS.BASE}/${deleteTarget.uuid}`;
    }
    return '';
  }, [deleteTarget]);

  const {
    mutate: removeGroup,
    loading: deleting,
    error: deleteError,
    reset: resetDelete,
  } = useMutation<unknown, void>(deleteEndpoint, 'DELETE');

  const {
    members,
    loading: membersLoading,
    error: membersError,
    mutating: membersMutating,
    addMember,
    updateMemberRole,
    removeMember,
    reset: resetMembers,
  } = useGroupMembers(membersTarget?.uuid ?? null);

  useEffect(() => {
    if (!membersTarget) {
      resetMembers();
    }
  }, [membersTarget, resetMembers]);

  const handleOpenGroupDocuments = (groupUUID: string) => {
    navigate({
      pathname: ROUTES.DOCUMENTS,
      search: `?group=${groupUUID}`,
    });
  };

  const handleOpenCreate = () => {
    resetCreate();
    setFormState({ mode: 'create' });
  };

  const handleOpenEdit = (group: GroupItemType) => {
    resetUpdate();
    setFormState({ mode: 'edit', group });
  };

  const handleCloseForm = () => {
    setFormState(null);
    resetCreate();
    resetUpdate();
  };

  const handleOpenMembers = (group: GroupItemType) => {
    if (group.role !== GROUP_ROLES.AUTHOR) {
      setSnackbar({
        message: t('groups.membersForbidden'),
        severity: 'warning',
      });
      return;
    }
    setMembersTarget(group);
  };

  const handleCloseMembers = () => {
    setMembersTarget(null);
  };

  const handleDeleteRequest = (group: GroupItemType) => {
    if (deleting && deleteTarget?.uuid === group.uuid) {
      return;
    }
    resetDelete();
    setDeleteTarget(group);
  };

  const handleCloseDelete = () => {
    setDeleteTarget(null);
    resetDelete();
  };

  const handleSubmitGroup = async (name: string) => {
    if (!formState) {
      return;
    }

    try {
      if (formState.mode === 'create') {
        await createGroup({ name });
        setSnackbar({
          message: t('groups.createSuccess'),
          severity: 'success',
        });
      } else if (formState.mode === 'edit' && formState.group) {
        await updateGroup({ name });
        setSnackbar({
          message: t('groups.updateSuccess'),
          severity: 'success',
        });
      }

      await refetch();
      handleCloseForm();
    } catch (mutationError) {
      const fallback =
        formState.mode === 'create'
          ? t('groups.createError')
          : t('groups.updateError');

      setSnackbar({
        message: fallback,
        severity: 'error',
      });
      console.error('Group mutation failed', mutationError);
    }
  };

  const handleConfirmDelete = async () => {
    if (!deleteTarget) {
      return;
    }

    try {
      await removeGroup();
      setSnackbar({
        message: t('groups.deleteSuccess'),
        severity: 'success',
      });
      await refetch();
      handleCloseDelete();
    } catch (mutationError) {
      setSnackbar({
        message: t('groups.deleteError'),
        severity: 'error',
      });
      console.error('Group deletion failed', mutationError);
    }
  };

  const handleSnackbarClose = () => {
    setSnackbar(null);
  };

  const extractErrorMessage = (error: unknown, fallback: string) => {
    if (error instanceof Error && error.message) {
      return error.message;
    }
    return fallback;
  };

  const handleAddMember = async (userUUID: string, role: GroupRole) => {
    try {
      await addMember(userUUID, role);
      setSnackbar({
        message: t('groups.membersAddSuccess'),
        severity: 'success',
      });
    } catch (error) {
      const message = extractErrorMessage(error, t('groups.membersAddError'));
      setSnackbar({ message, severity: 'error' });
      throw error instanceof Error ? error : new Error(message);
    }
  };

  const handleUpdateMemberRole = async (userUUID: string, role: GroupRole) => {
    try {
      await updateMemberRole(userUUID, role);
      setSnackbar({
        message: t('groups.membersUpdateSuccess'),
        severity: 'success',
      });
    } catch (error) {
      const message = extractErrorMessage(
        error,
        t('groups.membersUpdateError'),
      );
      setSnackbar({ message, severity: 'error' });
      throw error instanceof Error ? error : new Error(message);
    }
  };

  const handleRemoveMember = async (userUUID: string) => {
    try {
      await removeMember(userUUID);
      setSnackbar({
        message: t('groups.membersRemoveSuccess'),
        severity: 'success',
      });
    } catch (error) {
      const message = extractErrorMessage(
        error,
        t('groups.membersRemoveError'),
      );
      setSnackbar({ message, severity: 'error' });
      throw error instanceof Error ? error : new Error(message);
    }
  };

  const isFormOpen = Boolean(formState);
  const isCreateMode = formState?.mode === 'create';
  const formGroup = formState?.group;

  const formLoading = isCreateMode ? creating : updating;
  const formErrorMessage = isCreateMode
    ? createError
      ? createError.message || t('groups.createError')
      : null
    : updateError
      ? updateError.message || t('groups.updateError')
      : null;

  const deleteErrorMessage = deleteError
    ? deleteError.message || t('groups.deleteError')
    : null;

  const deleteDescriptionTemplate = t('groups.deleteConfirmDescription');
  const deleteDescription = deleteTarget
    ? deleteDescriptionTemplate.replace('{name}', deleteTarget.name)
    : deleteDescriptionTemplate.replace('{name}', '');

  const roleNames = useMemo(
    () => ({
      [GROUP_ROLES.AUTHOR]: t('groups.role.author'),
      [GROUP_ROLES.COAUTHOR]: t('groups.role.coauthor'),
      [GROUP_ROLES.REVIEWER]: t('groups.role.reviewer'),
    }),
    [t],
  );

  const canManageMembers = (group: GroupItemType) =>
    group.role === GROUP_ROLES.AUTHOR;
  const canEditGroup = (group: GroupItemType) =>
    group.role === GROUP_ROLES.AUTHOR;
  const canDeleteGroup = (group: GroupItemType) =>
    group.role === GROUP_ROLES.AUTHOR;

  return (
    <CenteredContent className={className}>
      <PageCard>
        <Stack spacing={3}>
          <PageHeader
            title={t('groups.title')}
            subtitle={t('groups.subtitle')}
          />

          <Button
            variant='contained'
            startIcon={<AddIcon />}
            onClick={handleOpenCreate}
            sx={{ alignSelf: { xs: 'stretch', sm: 'flex-end' } }}
          >
            {t('groups.createButton')}
          </Button>

          {loading && <LoadingSpinner />}

          {error && (
            <ErrorAlert
              message={t('groups.error')}
              onRetry={refetch}
              retryText={t('groups.refresh')}
            />
          )}

          {!loading && !error && groups.length === 0 && (
            <Typography color='text.secondary' align='center'>
              {t('groups.empty')}
            </Typography>
          )}

          {groups.length > 0 && (
            <GroupsList
              groups={groups}
              onGroupClick={handleOpenGroupDocuments}
              onGroupEdit={handleOpenEdit}
              onGroupDelete={handleDeleteRequest}
              totalLabel={t('groups.totalLabel')}
              createdAtLabel={t('groups.createdAt')}
              editLabel={t('groups.editLabel')}
              deleteLabel={t('groups.deleteLabel')}
              onGroupManageMembers={handleOpenMembers}
              manageMembersLabel={t('groups.manageMembersButton')}
              roleLabel={t('groups.roleLabel')}
              roleNames={roleNames}
              canManageMembers={canManageMembers}
              canEditGroup={canEditGroup}
              canDeleteGroup={canDeleteGroup}
            />
          )}
        </Stack>
      </PageCard>

      <GroupFormDialog
        open={isFormOpen}
        title={
          isCreateMode
            ? t('groups.createDialogTitle')
            : t('groups.editDialogTitle')
        }
        confirmLabel={
          isCreateMode ? t('groups.createConfirm') : t('groups.updateConfirm')
        }
        cancelLabel={t('groups.cancel')}
        nameLabel={t('groups.nameLabel')}
        namePlaceholder={t('groups.namePlaceholder')}
        nameHelperText={t('groups.nameHelper')}
        initialName={formGroup?.name}
        loading={formLoading}
        errorMessage={formErrorMessage}
        onClose={handleCloseForm}
        onSubmit={handleSubmitGroup}
      />

      <ConfirmDialog
        open={Boolean(deleteTarget)}
        title={t('groups.deleteConfirmTitle')}
        description={deleteDescription}
        confirmLabel={t('groups.deleteConfirmAccept')}
        cancelLabel={t('groups.cancel')}
        loading={deleting}
        errorMessage={deleteErrorMessage}
        onCancel={handleCloseDelete}
        onConfirm={handleConfirmDelete}
      />

      <Snackbar
        open={Boolean(snackbar)}
        autoHideDuration={4000}
        onClose={handleSnackbarClose}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        {snackbar ? (
          <Alert
            onClose={handleSnackbarClose}
            severity={snackbar.severity}
            sx={{ width: '100%' }}
          >
            {snackbar.message}
          </Alert>
        ) : undefined}
      </Snackbar>

      <GroupMembersDialog
        open={Boolean(membersTarget)}
        group={membersTarget}
        members={membersTarget ? members : []}
        loading={membersLoading && Boolean(membersTarget)}
        mutating={membersMutating}
        error={membersError}
        onClose={handleCloseMembers}
        onAddMember={handleAddMember}
        onUpdateMemberRole={handleUpdateMemberRole}
        onRemoveMember={handleRemoveMember}
      />
    </CenteredContent>
  );
};

export default Groups;
