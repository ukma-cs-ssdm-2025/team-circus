import { useEffect, useMemo, useState } from 'react';
import { Stack, Typography } from '@mui/material';
import { useLocation } from 'react-router-dom';
import { CenteredContent, PageCard, PageHeader, ErrorAlert, LoadingSpinner, DocumentFilters, DocumentsGrid } from '../components';
import { useLanguage } from '../contexts/LanguageContext';
import { useApi } from '../hooks';
import { API_ENDPOINTS } from '../constants';
import type { BaseComponentProps, DocumentsResponse, GroupsResponse, DocumentFilters as DocumentFiltersType, GroupOption } from '../types';

type DocumentsProps = BaseComponentProps;

const Documents = ({ className = '' }: DocumentsProps) => {
  const { t } = useLanguage();
  const location = useLocation();
  const [filters, setFilters] = useState<DocumentFiltersType>({
    selectedGroup: 'all',
    searchTerm: '',
  });

  const {
    data: documentsData,
    loading: documentsLoading,
    error: documentsError,
    refetch: refetchDocuments,
  } = useApi<DocumentsResponse>(API_ENDPOINTS.DOCUMENTS.BASE);

  const {
    data: groupsData,
    loading: groupsLoading,
  } = useApi<GroupsResponse>(API_ENDPOINTS.GROUPS.BASE);

  const documents = useMemo(() => documentsData?.documents ?? [], [documentsData]);
  const groups = useMemo(() => groupsData?.groups ?? [], [groupsData]);

  const groupOptions = useMemo((): GroupOption[] => {
    return groups.map(group => ({
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
    return documents.filter(document => {
      const matchesGroup = filters.selectedGroup === 'all' || document.group_uuid === filters.selectedGroup;
      const matchesSearch = document.name.toLowerCase().includes(filters.searchTerm.toLowerCase());
      return matchesGroup && matchesSearch;
    });
  }, [documents, filters.selectedGroup, filters.searchTerm]);

  const handleGroupChange = (value: string) => {
    setFilters(prev => ({ ...prev, selectedGroup: value }));
  };

  const handleSearchChange = (value: string) => {
    setFilters(prev => ({ ...prev, searchTerm: value }));
  };

  const isLoading = documentsLoading || groupsLoading;

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const groupFromQuery = params.get('group');

    if (!groupFromQuery) {
      return;
    }

    const isKnownGroup = groups.some(group => group.uuid === groupFromQuery);
    if (isKnownGroup) {
      setFilters(prev => ({ ...prev, selectedGroup: groupFromQuery }));
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
            <Typography color="text.secondary" align="center">
              {t('documents.empty')}
            </Typography>
          )}

          <DocumentsGrid
            documents={filteredDocuments}
            groupNameByUUID={groupNameByUUID}
            createdAtLabel={t('documents.createdAt')}
            noContentLabel={t('documents.noContent')}
            groupUnknownLabel={t('documents.groupUnknown')}
          />
        </Stack>
      </PageCard>
    </CenteredContent>
  );
};

export default Documents;
