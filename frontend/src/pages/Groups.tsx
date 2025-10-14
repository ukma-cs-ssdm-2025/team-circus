import { Typography } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { CenteredContent, PageCard, PageHeader, ErrorAlert, LoadingSpinner, GroupsList } from '../components';
import { useLanguage } from '../contexts/LanguageContext';
import { useApi } from '../hooks';
import { API_ENDPOINTS, ROUTES } from '../constants';
import type { BaseComponentProps, GroupsResponse } from '../types';

type GroupsProps = BaseComponentProps;

const Groups = ({ className = '' }: GroupsProps) => {
  const { t } = useLanguage();
  const navigate = useNavigate();
  const { data, loading, error, refetch } = useApi<GroupsResponse>(API_ENDPOINTS.GROUPS.BASE);
  const groups = data?.groups ?? [];

  const handleOpenGroupDocuments = (groupUUID: string) => {
    navigate({
      pathname: ROUTES.DOCUMENTS,
      search: `?group=${groupUUID}`,
    });
  };

  return (
    <CenteredContent className={className}>
      <PageCard>
        <PageHeader
          title={t('groups.title')}
          subtitle={t('groups.subtitle')}
        />

        {loading && <LoadingSpinner />}

        {error && (
          <ErrorAlert
            message={t('groups.error')}
            onRetry={refetch}
            retryText={t('groups.refresh')}
          />
        )}

        {!loading && !error && groups.length === 0 && (
          <Typography color="text.secondary" align="center">
            {t('groups.empty')}
          </Typography>
        )}

        {groups.length > 0 && (
          <GroupsList
            groups={groups}
            onGroupClick={handleOpenGroupDocuments}
            totalLabel={t('groups.totalLabel')}
            createdAtLabel={t('groups.createdAt')}
          />
        )}
      </PageCard>
    </CenteredContent>
  );
};

export default Groups;
