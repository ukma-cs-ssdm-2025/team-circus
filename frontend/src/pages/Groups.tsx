import {
  Alert,
  Box,
  Button,
  CircularProgress,
  Divider,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  Stack,
  Typography,
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { CenteredContent, PageCard, PageHeader } from '../components/common';
import { useLanguage } from '../contexts/LanguageContext';
import { useApi } from '../hooks';
import { API_ENDPOINTS, ROUTES } from '../constants';
import { formatDate } from '../utils';
import type { BaseComponentProps } from '../types';

interface GroupsResponse {
  groups: GroupItem[];
}

interface GroupItem {
  uuid: string;
  name: string;
  created_at: string;
}

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

        {loading && (
          <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
            <CircularProgress />
          </Box>
        )}

        {error && (
          <Alert
            severity="error"
            sx={{ mb: 3, display: 'flex', alignItems: 'center', gap: 2 }}
          >
            <Box component="span">{t('groups.error')}</Box>
            <Button variant="outlined" size="small" onClick={refetch}>
              {t('groups.refresh')}
            </Button>
          </Alert>
        )}

        {!loading && !error && groups.length === 0 && (
          <Typography color="text.secondary" align="center">
            {t('groups.empty')}
          </Typography>
        )}

        {groups.length > 0 && (
          <Stack spacing={2}>
            <Typography variant="body2" color="text.secondary">
              {`${t('groups.totalLabel')} ${groups.length}`}
            </Typography>

            <List disablePadding>
              {groups.map((group, index) => (
                <Box key={group.uuid}>
                  <ListItem disablePadding alignItems="flex-start">
                    <ListItemButton
                      sx={{ borderRadius: 2 }}
                      onClick={() => handleOpenGroupDocuments(group.uuid)}
                    >
                      <ListItemText
                        primary={
                          <Typography variant="h6" sx={{ fontWeight: 600 }}>
                            {group.name}
                          </Typography>
                        }
                        secondary={
                          <Typography variant="body2" color="text.secondary">
                            {`${t('groups.createdAt')}: ${formatDate(group.created_at)}`}
                          </Typography>
                        }
                      />
                    </ListItemButton>
                  </ListItem>
                  {index < groups.length - 1 && <Divider sx={{ my: 1 }} />}
                </Box>
              ))}
            </List>
          </Stack>
        )}
      </PageCard>
    </CenteredContent>
  );
};

export default Groups;
