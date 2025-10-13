import { useEffect, useMemo, useState } from 'react';
import {
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  CircularProgress,
  FormControl,
  Grid,
  InputLabel,
  MenuItem,
  Select,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import { useLocation } from 'react-router-dom';
import { CenteredContent, PageCard, PageHeader } from '../components/common';
import { useLanguage } from '../contexts/LanguageContext';
import { useApi } from '../hooks';
import { API_ENDPOINTS } from '../constants';
import { formatDate, truncateText } from '../utils';
import type { BaseComponentProps } from '../types';

interface DocumentsResponse {
  documents: DocumentItem[];
}

interface GroupsResponse {
  groups: GroupItem[];
}

interface DocumentItem {
  uuid: string;
  name: string;
  content: string;
  group_uuid: string;
  created_at: string;
}

interface GroupItem {
  uuid: string;
  name: string;
}

type DocumentsProps = BaseComponentProps;

const Documents = ({ className = '' }: DocumentsProps) => {
  const { t } = useLanguage();
  const location = useLocation();
  const [selectedGroup, setSelectedGroup] = useState<string>('all');
  const [searchTerm, setSearchTerm] = useState<string>('');

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

  const groupOptions = useMemo(() => {
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
      const matchesGroup = selectedGroup === 'all' || document.group_uuid === selectedGroup;
      const matchesSearch = document.name.toLowerCase().includes(searchTerm.toLowerCase());
      return matchesGroup && matchesSearch;
    });
  }, [documents, selectedGroup, searchTerm]);

  const handleGroupChange = (value: string) => {
    setSelectedGroup(value);
  };

  const handleSearchChange = (value: string) => {
    setSearchTerm(value);
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
      setSelectedGroup(groupFromQuery);
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
          <Stack direction={{ xs: 'column', md: 'row' }} spacing={2}>
            <FormControl fullWidth>
              <InputLabel>{t('documents.filterGroup')}</InputLabel>
              <Select
                value={selectedGroup}
                label={t('documents.filterGroup')}
                onChange={(event) => handleGroupChange(event.target.value)}
              >
                <MenuItem value="all">{t('documents.filterAll')}</MenuItem>
                {groupOptions.map(option => (
                  <MenuItem key={option.value} value={option.value}>
                    {option.label}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>

            <TextField
              fullWidth
              value={searchTerm}
              label={t('documents.searchPlaceholder')}
              onChange={(event) => handleSearchChange(event.target.value)}
            />
          </Stack>

          {isLoading && (
            <Box sx={{ display: 'flex', justifyContent: 'center', py: 6 }}>
              <CircularProgress />
            </Box>
          )}

          {documentsError && (
            <Alert
              severity="error"
              sx={{ display: 'flex', alignItems: 'center', gap: 2 }}
            >
              <Box component="span">{t('documents.error')}</Box>
              <Button variant="outlined" size="small" onClick={refetchDocuments}>
                {t('documents.refresh')}
              </Button>
            </Alert>
          )}

          {!isLoading && !documentsError && filteredDocuments.length === 0 && (
            <Typography color="text.secondary" align="center">
              {t('documents.empty')}
            </Typography>
          )}

          <Grid container spacing={3}>
            {filteredDocuments.map(document => (
              <Grid item xs={12} md={6} key={document.uuid}>
                <Card
                  sx={{
                    height: '100%',
                    display: 'flex',
                    flexDirection: 'column',
                    borderRadius: 3,
                    boxShadow: '0 10px 30px rgba(0, 0, 0, 0.08)',
                  }}
                >
                  <CardContent sx={{ flex: 1 }}>
                    <Stack spacing={2}>
                      <Stack direction="row" spacing={1} alignItems="center">
                        <Typography variant="h6" sx={{ fontWeight: 600 }}>
                          {document.name}
                        </Typography>
                        <Chip
                          label={groupNameByUUID[document.group_uuid] || t('documents.groupUnknown')}
                          size="small"
                          color="primary"
                          sx={{ fontWeight: 600 }}
                        />
                      </Stack>

                      <Typography variant="body2" color="text.secondary">
                        {truncateText(document.content || t('documents.noContent'), 180)}
                      </Typography>

                      <Typography variant="body2" color="text.secondary">
                        {`${t('documents.createdAt')}: ${formatDate(document.created_at)}`}
                      </Typography>
                    </Stack>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>
        </Stack>
      </PageCard>
    </CenteredContent>
  );
};

export default Documents;
