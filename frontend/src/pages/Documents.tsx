import { useEffect, useMemo, useState } from 'react';
import { Alert, Button, Snackbar, Stack, Typography } from '@mui/material';
import type { AlertColor } from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import { useLocation, useNavigate } from 'react-router-dom';
import {
  CenteredContent,
  PageCard,
  PageHeader,
  ErrorAlert,
  LoadingSpinner,
  DocumentFilters,
  DocumentsGrid,
  DocumentFormDialog,
  ConfirmDialog,
} from '../components';
import { useLanguage } from '../contexts/LanguageContext';
import { useApi, useMutation } from '../hooks';
import { API_ENDPOINTS, ROUTES } from '../constants';
import type {
  BaseComponentProps,
  DocumentsResponse,
  GroupsResponse,
  DocumentItem,
  DocumentFilters as DocumentFiltersType,
  GroupOption,
  CreateDocumentPayload,
} from '../types';

type DocumentsProps = BaseComponentProps;

const Documents = ({ className = '' }: DocumentsProps) => {
  const { t } = useLanguage();
  const location = useLocation();
  const navigate = useNavigate();
  const [filters, setFilters] = useState<DocumentFiltersType>({
    selectedGroup: 'all',
    searchTerm: '',
  });
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [deleteTarget, setDeleteTarget] = useState<DocumentItem | null>(null);
  const [snackbar, setSnackbar] = useState<{
    message: string;
    severity: AlertColor;
  } | null>(null);

  const {
    data: documentsData,
    loading: documentsLoading,
    error: documentsError,
    refetch: refetchDocuments,
  } = useApi<DocumentsResponse>(API_ENDPOINTS.DOCUMENTS.BASE);

  const { data: groupsData, loading: groupsLoading } = useApi<GroupsResponse>(
    API_ENDPOINTS.GROUPS.BASE,
  );

  const {
    mutate: createDocument,
    loading: creatingDocument,
    error: createError,
    reset: resetCreate,
  } = useMutation<DocumentItem, CreateDocumentPayload>(
    API_ENDPOINTS.DOCUMENTS.BASE,
    'POST',
  );

  const deleteEndpoint = useMemo(() => {
    if (deleteTarget) {
      return `${API_ENDPOINTS.DOCUMENTS.BASE}/${deleteTarget.uuid}`;
    }
    return '';
  }, [deleteTarget]);

  const {
    mutate: removeDocument,
    loading: deletingDocument,
    error: deleteError,
    reset: resetDelete,
  } = useMutation<unknown, void>(deleteEndpoint, 'DELETE');

  const documents = useMemo(
    () => documentsData?.documents ?? [],
    [documentsData],
  );
  const groups = useMemo(() => groupsData?.groups ?? [], [groupsData]);

  const groupOptions = useMemo((): GroupOption[] => {
    return groups.map((group) => ({
      value: group.uuid,
      label: group.name,
    }));
  }, [groups]);

  const groupNameByUUID = useMemo(() => {
    return groups.reduce<Record<string, string>>((acc, group) => {
      acc[group.uuid] = group.name;
      return acc;
    }, {});
  }, [groups]);

  const filteredDocuments = useMemo(() => {
    return documents.filter((document) => {
      const matchesGroup =
        filters.selectedGroup === 'all' ||
        document.group_uuid === filters.selectedGroup;
      const matchesSearch = document.name
        .toLowerCase()
        .includes(filters.searchTerm.toLowerCase());
      return matchesGroup && matchesSearch;
    });
  }, [documents, filters.selectedGroup, filters.searchTerm]);

  const handleGroupChange = (value: string) => {
    setFilters((prev) => ({ ...prev, selectedGroup: value }));
  };

  const handleSearchChange = (value: string) => {
    setFilters((prev) => ({ ...prev, searchTerm: value }));
  };

  const isLoading = documentsLoading || groupsLoading;

  const handleOpenCreate = () => {
    resetCreate();
    setIsCreateOpen(true);
  };

  const handleCloseCreate = () => {
    setIsCreateOpen(false);
    resetCreate();
  };

  const handleCreateDocument = async ({
    name,
    groupUUID,
  }: {
    name: string;
    groupUUID: string;
  }) => {
    try {
      const defaultContent = t('documents.defaultContent');
      const payload: CreateDocumentPayload = {
        name,
        group_uuid: groupUUID,
        content:
          defaultContent && defaultContent !== 'documents.defaultContent'
            ? defaultContent
            : ' ',
      };

      const created = await createDocument({
        ...payload,
      });

      setSnackbar({
        message: t('documents.createSuccess'),
        severity: 'success',
      });

      await refetchDocuments();
      handleCloseCreate();

      if (created?.uuid) {
        navigate(`${ROUTES.DOCUMENTS}/${created.uuid}`);
      }
    } catch (mutationError) {
      setSnackbar({
        message: t('documents.createError'),
        severity: 'error',
      });
      console.error('Document creation failed', mutationError);
    }
  };

  const handleRequestDelete = (document: DocumentItem) => {
    if (deletingDocument && deleteTarget?.uuid === document.uuid) {
      return;
    }
    resetDelete();
    setDeleteTarget(document);
  };

  const handleCloseDelete = () => {
    setDeleteTarget(null);
    resetDelete();
  };

  const handleConfirmDelete = async () => {
    if (!deleteTarget) {
      return;
    }

    try {
      await removeDocument();
      setSnackbar({
        message: t('documents.deleteSuccess'),
        severity: 'success',
      });
      await refetchDocuments();
      handleCloseDelete();
    } catch (mutationError) {
      setSnackbar({
        message: t('documents.deleteError'),
        severity: 'error',
      });
      console.error('Document deletion failed', mutationError);
    }
  };

  const handleSnackbarClose = () => {
    setSnackbar(null);
  };

  const createErrorMessage = createError
    ? createError.message || t('documents.createError')
    : null;
  const deleteErrorMessage = deleteError
    ? deleteError.message || t('documents.deleteError')
    : null;

  const deleteDescriptionTemplate = t('documents.deleteConfirmDescription');
  const deleteDescription = deleteTarget
    ? deleteDescriptionTemplate.replace('{name}', deleteTarget.name)
    : deleteDescriptionTemplate.replace('{name}', '');

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const groupFromQuery = params.get('group');

    if (!groupFromQuery) {
      return;
    }

    const isKnownGroup = groups.some((group) => group.uuid === groupFromQuery);
    if (isKnownGroup) {
      setFilters((prev) => ({ ...prev, selectedGroup: groupFromQuery }));
    }
  }, [location.search, groups]);

  return (
    <CenteredContent className={className}>
      <PageCard>
        <PageHeader
          title={t('documents.title')}
          subtitle={t('documents.subtitle')}
        />

        <Stack spacing={3}>
          <Stack
            direction={{ xs: 'column', sm: 'row' }}
            spacing={2}
            alignItems={{ xs: 'stretch', sm: 'center' }}
          >
            <Button
              variant='contained'
              startIcon={<AddIcon />}
              onClick={handleOpenCreate}
              disabled={groups.length === 0}
              sx={{ alignSelf: { xs: 'stretch', sm: 'flex-start' } }}
            >
              {t('documents.createButton')}
            </Button>
            {groups.length === 0 && !groupsLoading && (
              <Alert severity='info' sx={{ flex: 1 }}>
                {t('documents.noGroupsHint')}
              </Alert>
            )}
          </Stack>

          <DocumentFilters
            filters={filters}
            groupOptions={groupOptions}
            onGroupChange={handleGroupChange}
            onSearchChange={handleSearchChange}
            filterGroupLabel={t('documents.filterGroup')}
            filterAllLabel={t('documents.filterAll')}
            searchPlaceholder={t('documents.searchPlaceholder')}
          />

          {isLoading && <LoadingSpinner py={6} />}

          {documentsError && (
            <ErrorAlert
              message={t('documents.error')}
              onRetry={refetchDocuments}
              retryText={t('documents.refresh')}
            />
          )}

          {!isLoading && !documentsError && filteredDocuments.length === 0 && (
            <Typography color='text.secondary' align='center'>
              {t('documents.empty')}
            </Typography>
          )}

          <DocumentsGrid
            documents={filteredDocuments}
            groupNameByUUID={groupNameByUUID}
            createdAtLabel={t('documents.createdAt')}
            noContentLabel={t('documents.noContent')}
            groupUnknownLabel={t('documents.groupUnknown')}
            editLabel={t('documents.editLabel')}
            deleteLabel={t('documents.deleteLabel')}
            onDocumentDelete={handleRequestDelete}
          />
        </Stack>
      </PageCard>

      <DocumentFormDialog
        open={isCreateOpen}
        title={t('documents.createDialogTitle')}
        confirmLabel={t('documents.createConfirm')}
        cancelLabel={t('documents.cancel')}
        nameLabel={t('documents.nameLabel')}
        namePlaceholder={t('documents.namePlaceholder')}
        nameHelperText={t('documents.nameHelper')}
        groupLabel={t('documents.groupLabel')}
        groupOptions={groupOptions}
        loading={creatingDocument}
        errorMessage={createErrorMessage}
        onClose={handleCloseCreate}
        onSubmit={handleCreateDocument}
      />

      <ConfirmDialog
        open={Boolean(deleteTarget)}
        title={t('documents.deleteConfirmTitle')}
        description={deleteDescription}
        confirmLabel={t('documents.deleteConfirmAccept')}
        cancelLabel={t('documents.cancel')}
        loading={deletingDocument}
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
        {snackbar && (
          <Alert
            onClose={handleSnackbarClose}
            severity={snackbar.severity}
            sx={{ width: '100%' }}
          >
            {snackbar.message}
          </Alert>
        )}
      </Snackbar>
    </CenteredContent>
  );
};

export default Documents;
