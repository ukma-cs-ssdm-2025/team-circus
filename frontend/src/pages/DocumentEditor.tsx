import { useCallback, useEffect, useMemo, useRef, useState, type ChangeEvent } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Alert,
  Box,
  Button,
  Chip,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { ErrorAlert, LoadingSpinner } from '../components';
import { useLanguage } from '../contexts/LanguageContext';
import { useApi, useMutation } from '../hooks';
import { API_ENDPOINTS, ROUTES } from '../constants';
import { getWebSocketUrl } from '../config/env';
import type { BaseComponentProps, DocumentItem, GroupsResponse } from '../types';

type UpdateDocumentPayload = {
  name: string;
  content: string;
};

type DocumentSocketMessage = {
  type: 'init' | 'update' | 'error';
  document?: DocumentItem;
  error?: string;
};

type DocumentEditorProps = BaseComponentProps;

const markdownStyles = {
  fontSize: 15,
  lineHeight: 1.45,
  '& > *:first-of-type': { marginTop: 0 },
  '& > *:last-child': { marginBottom: 0 },
  '& h1': { fontSize: 28, fontWeight: 700, margin: '12px 0 6px' },
  '& h2': { fontSize: 22, fontWeight: 700, margin: '10px 0 4px' },
  '& h3': { fontSize: 19, fontWeight: 600, margin: '8px 0 3px' },
  '& p': { margin: '0 0 4px' },
  '& ul, & ol': {
    paddingLeft: '1.25ch',
    margin: '0 0 6px',
    listStyle: 'none',
  },
  '& li': {
    margin: '0 0 2px',
    paddingLeft: 0,
    display: 'flex',
    alignItems: 'flex-start',
    gap: '1.3ch',
  },
  '& ul li::before': {
    content: '"•"',
    display: 'inline-block',
    width: '1.25ch',
    color: 'inherit',
  },
  '& ol': { counterReset: 'markdown-ol' },
  '& ol li': { counterIncrement: 'markdown-ol' },
  '& ol li::before': {
    content: 'counter(markdown-ol) "."',
    display: 'inline-block',
    width: '2.8ch',
    color: 'inherit',
    fontVariantNumeric: 'tabular-nums',
  },
  '& li:last-child': { marginBottom: 0 },
  '& li > p': { margin: 0, display: 'inline' },
  '& blockquote': {
    borderLeft: '4px solid rgba(0, 0, 0, 0.1)',
    paddingLeft: 12,
    color: 'text.secondary',
    fontStyle: 'italic',
    margin: '6px 0',
  },
  '& code': {
    fontFamily: 'ui-monospace, SFMono-Regular, SFMono, Menlo, Monaco, Consolas, Liberation Mono, Courier New, monospace',
    backgroundColor: 'rgba(15, 23, 42, 0.06)',
    borderRadius: 1,
    padding: '2px 4px',
    fontSize: 14,
  },
  '& pre': {
    fontFamily: 'ui-monospace, SFMono-Regular, SFMono, Menlo, Monaco, Consolas, Liberation Mono, Courier New, monospace',
    backgroundColor: 'rgba(15, 23, 42, 0.08)',
    borderRadius: 2,
    padding: 12,
    overflowX: 'auto',
    margin: '8px 0',
  },
  '& table': {
    width: '100%',
    borderCollapse: 'collapse',
    margin: '8px 0',
  },
  '& th, & td': {
    border: '1px solid rgba(148, 163, 184, 0.4)',
    padding: 10,
    textAlign: 'left',
  },
} as const;

const DocumentEditor = ({ className = '' }: DocumentEditorProps) => {
  const { t } = useLanguage();
  const navigate = useNavigate();
  const { uuid } = useParams<{ uuid: string }>();
  const documentId = uuid ?? '';
  const endpoint = useMemo(
    () => (documentId ? `${API_ENDPOINTS.DOCUMENTS.BASE}/${documentId}` : ''),
    [documentId]
  );

  const {
    data: documentData,
    loading: documentLoading,
    error: documentError,
    refetch: refetchDocument,
    mutate: updateDocumentCache,
  } = useApi<DocumentItem>(endpoint, {
    immediate: Boolean(documentId),
  });

  const {
    data: groupsData,
  } = useApi<GroupsResponse>(API_ENDPOINTS.GROUPS.BASE);

  const { mutate: updateDocument, loading: httpSaving } = useMutation<DocumentItem, UpdateDocumentPayload>(
    endpoint,
    'PUT'
  );

  const [name, setName] = useState('');
  const [content, setContent] = useState('');
  const [lastSaved, setLastSaved] = useState<UpdateDocumentPayload | null>(null);
  const [saveStatus, setSaveStatus] = useState<'idle' | 'success' | 'error'>('idle');
  const [pendingWsSave, setPendingWsSave] = useState(false);
  const [wsStatus, setWsStatus] = useState<'disconnected' | 'connecting' | 'connected' | 'error'>('disconnected');
  const [wsError, setWsError] = useState<string | null>(null);
  const socketRef = useRef<WebSocket | null>(null);
  const pendingWsSaveRef = useRef(false);
  const pendingWsOriginRef = useRef<'manual' | 'auto' | null>(null);

  useEffect(() => {
    if (documentData) {
      setName(documentData.name ?? '');
      setContent(documentData.content ?? '');
      setLastSaved({
        name: documentData.name ?? '',
        content: documentData.content ?? '',
      });
      setSaveStatus('idle');
      setPendingWsSave(false);
      pendingWsSaveRef.current = false;
      setWsError(null);
    }
  }, [documentData]);

  useEffect(() => {
    pendingWsSaveRef.current = pendingWsSave;
  }, [pendingWsSave]);

  const documentGroupUUID = documentData?.group_uuid ?? '';
  const groupName = useMemo(() => {
    if (!documentGroupUUID || !groupsData?.groups) {
      return '';
    }

    const group = groupsData.groups.find(candidate => candidate.uuid === documentGroupUUID);
    return group?.name ?? '';
  }, [documentGroupUUID, groupsData?.groups]);

  const isNameValid = name.trim().length > 0;
  const isDirty = lastSaved !== null && (lastSaved.name !== name || lastSaved.content !== content);
  const isSaving = httpSaving || pendingWsSave;
  const canUseWebSocket = wsStatus === 'connected' && socketRef.current?.readyState === WebSocket.OPEN;
  const isSaveDisabled = !isNameValid || !isDirty || isSaving || !documentId;

  const sendWebsocketUpdate = useCallback((nextName: string, nextContent: string, origin: 'manual' | 'auto' = 'manual') => {
    if (!canUseWebSocket) {
      return false;
    }

    const socket = socketRef.current;
    if (!socket || socket.readyState !== WebSocket.OPEN) {
      return false;
    }

    try {
      setPendingWsSave(true);
      pendingWsSaveRef.current = true;
      pendingWsOriginRef.current = origin;
      setSaveStatus('idle');
      if (origin === 'manual') {
        setWsError(null);
      }

      socket.send(JSON.stringify({
        type: 'update',
        name: nextName,
        content: nextContent,
      }));
      return true;
    } catch (error) {
      console.error('Failed to send websocket update', error);
      setPendingWsSave(false);
      pendingWsSaveRef.current = false;
      pendingWsOriginRef.current = null;
      if (origin === 'manual') {
        setSaveStatus('error');
        setWsError('Unable to send live update');
      }
      return false;
    }
  }, [canUseWebSocket]);

  useEffect(() => {
    if (!documentId) {
      return;
    }

    const wsUrl = getWebSocketUrl(`${API_ENDPOINTS.DOCUMENTS.BASE}/${documentId}/ws`);
    setWsStatus('connecting');
    setWsError(null);

    const socket = new WebSocket(wsUrl);
    socketRef.current = socket;

    socket.onopen = () => {
      setWsStatus('connected');
    };

    socket.onmessage = event => {
      try {
        const message: DocumentSocketMessage = JSON.parse(event.data);
        if ((message.type === 'init' || message.type === 'update') && message.document) {
          const incoming = message.document;
          setName(incoming.name ?? '');
          setContent(incoming.content ?? '');
          setLastSaved({
            name: incoming.name ?? '',
            content: incoming.content ?? '',
          });
          updateDocumentCache(incoming);
          if (pendingWsSaveRef.current) {
            if (pendingWsOriginRef.current === 'manual') {
              setSaveStatus('success');
            } else {
              setSaveStatus(prev => (prev === 'error' ? 'idle' : prev));
            }
            setPendingWsSave(false);
            pendingWsSaveRef.current = false;
            pendingWsOriginRef.current = null;
            setWsError(null);
          } else {
            setSaveStatus(prev => (prev === 'error' ? 'idle' : prev));
          }
        } else if (message.type === 'error') {
          setWsError(message.error ?? 'Live update error');
          setSaveStatus(pendingWsOriginRef.current === 'manual' ? 'error' : 'idle');
          setPendingWsSave(false);
          pendingWsSaveRef.current = false;
          pendingWsOriginRef.current = null;
        }
      } catch (error) {
        console.error('Failed to parse websocket message', error);
      }
    };

    socket.onerror = () => {
      setWsStatus('error');
      setWsError('Live updates connection failed');
      setSaveStatus('error');
      setPendingWsSave(false);
      pendingWsSaveRef.current = false;
      pendingWsOriginRef.current = null;
    };

    socket.onclose = () => {
      setWsStatus('disconnected');
      setPendingWsSave(false);
      pendingWsSaveRef.current = false;
      pendingWsOriginRef.current = null;
    };

    return () => {
      socket.close();
      socketRef.current = null;
      setWsStatus('disconnected');
      setPendingWsSave(false);
      pendingWsSaveRef.current = false;
      pendingWsOriginRef.current = null;
    };
  }, [documentId, updateDocumentCache]);

  const handleNameChange = (event: ChangeEvent<HTMLInputElement>) => {
    if (saveStatus !== 'idle' || wsError) {
      setSaveStatus('idle');
      setWsError(null);
    }
    setName(event.target.value);
  };

  const handleContentChange = (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    if (saveStatus !== 'idle' || wsError) {
      setSaveStatus('idle');
      setWsError(null);
    }
    setContent(event.target.value);
  };

  const handleSave = async () => {
    if (isSaveDisabled || !documentId) {
      return;
    }

    const trimmedName = name.trim();

    if (sendWebsocketUpdate(trimmedName, content, 'manual')) {
      return;
    }

    try {
      const updated = await updateDocument({
        name: trimmedName,
        content,
      });

      updateDocumentCache(updated);
      setLastSaved({
        name: updated.name ?? '',
        content: updated.content ?? '',
      });
      setSaveStatus('success');
      setWsError(null);
    } catch (error) {
      console.error('Failed to save document', error);
      setSaveStatus('error');
      setWsError('Unable to save changes');
      pendingWsOriginRef.current = null;
    }
  };

  useEffect(() => {
    if (!documentId || !canUseWebSocket) {
      return;
    }
    if (!isDirty || pendingWsSaveRef.current || !lastSaved) {
      return;
    }

    const trimmedName = name.trim();
    const timer = window.setTimeout(() => {
      sendWebsocketUpdate(trimmedName, content, 'auto');
    }, 600);

    return () => {
      window.clearTimeout(timer);
    };
  }, [documentId, canUseWebSocket, isDirty, name, content, sendWebsocketUpdate, lastSaved]);

  if (!documentId) {
    return (
      <Box className={className} sx={{ px: { xs: 2, md: 4 }, py: { xs: 3, md: 4 } }}>
        <Stack spacing={3}>
          <Typography variant="h5" fontWeight={700}>
            {t('documentEditor.fallbackTitle')}
          </Typography>
          <ErrorAlert message={t('documentEditor.notFound')} />
          <Box>
            <Button variant="outlined" onClick={() => navigate(ROUTES.DOCUMENTS)}>
              {t('documentEditor.backToList')}
            </Button>
          </Box>
        </Stack>
      </Box>
    );
  }

  return (
    <Box
      className={className}
      sx={{
        display: 'flex',
        flexDirection: 'column',
        minHeight: '100vh',
        px: { xs: 2, md: 4 },
        py: { xs: 3, md: 4 },
        gap: 3,
      }}
    >
      <Stack
        direction={{ xs: 'column', sm: 'row' }}
        spacing={2}
        justifyContent="space-between"
        alignItems={{ xs: 'flex-start', sm: 'center' }}
      >
        <Stack direction="row" spacing={2} alignItems="center">
          <Button variant="outlined" onClick={() => navigate(ROUTES.DOCUMENTS)}>
            {t('documentEditor.backToList')}
          </Button>
          <Chip
            label={groupName || t('documents.groupUnknown')}
            color="primary"
            sx={{ fontWeight: 600 }}
          />
        </Stack>
        <Button variant="contained" onClick={handleSave} disabled={isSaveDisabled}>
          {isSaving ? t('documentEditor.savingButton') : t('documentEditor.saveButton')}
        </Button>
      </Stack>

      {documentError && (
        <ErrorAlert
          message={t('documentEditor.loadError')}
          onRetry={refetchDocument}
          retryText={t('documents.refresh')}
        />
      )}

      {saveStatus !== 'idle' && (
        <Alert
          severity={saveStatus === 'success' ? 'success' : 'error'}
          onClose={() => {
            setSaveStatus('idle');
            setWsError(null);
          }}
        >
          {saveStatus === 'success'
            ? t('documentEditor.saveSuccess')
            : wsError ?? t('documentEditor.saveError')}
        </Alert>
      )}

      {documentLoading && <LoadingSpinner py={6} />}

      {!documentLoading && !documentError && !documentData && (
        <Typography color="text.secondary">{t('documentEditor.notFound')}</Typography>
      )}

      {!documentLoading && documentData && (
        <Box
          sx={{
            flex: 1,
            display: 'flex',
            flexDirection: { xs: 'column', md: 'row' },
            gap: { xs: 3, md: 4 },
            alignItems: 'stretch',
          }}
        >
          <Box
            sx={{
              flex: 1,
              display: 'flex',
              flexDirection: 'column',
              gap: 2,
              minHeight: { xs: '32vh', md: '55vh' },
              height: { md: '70vh' },
              backgroundColor: 'rgba(148, 163, 184, 0.08)',
              borderRadius: 2,
              border: '1px solid rgba(148, 163, 184, 0.24)',
              padding: { xs: 2, md: 3 },
              boxShadow: '0 10px 30px rgba(15, 23, 42, 0.18)',
            }}
          >
            <TextField
              label={t('documentEditor.nameLabel')}
              placeholder={t('documentEditor.namePlaceholder')}
              value={name}
              onChange={handleNameChange}
              fullWidth
              error={!isNameValid}
              helperText={!isNameValid ? t('documentEditor.nameRequired') : undefined}
              FormHelperTextProps={{ sx: { minHeight: 0, mt: isNameValid ? 0 : 0.5 } }}
              sx={{
                '& .MuiOutlinedInput-root': {
                  borderRadius: 1.5,
                },
                '& .MuiInputBase-input': {
                  padding: '12px 14px',
                },
              }}
            />

            <TextField
              label={t('documentEditor.contentLabel')}
              placeholder={t('documentEditor.contentPlaceholder')}
              value={content}
              onChange={handleContentChange}
              fullWidth
              multiline
              minRows={16}
              sx={{
                flex: 1,
                '& .MuiOutlinedInput-root': {
                  alignItems: 'stretch',
                  borderRadius: 1.5,
                  height: '100%',
                  display: 'flex',
                },
                '& .MuiInputBase-input': {
                  fontFamily: 'ui-monospace, SFMono-Regular, SFMono, Menlo, Monaco, Consolas, Liberation Mono, Courier New, monospace',
                  padding: '12px 14px',
                },
                '& textarea': {
                  flex: 1,
                  overflowY: 'auto',
                  resize: 'none',
                  height: '100% !important',
                  maxHeight: '100% !important',
                  minHeight: '0 !important',
                },
              }}
            />
          </Box>

          <Box
            sx={{
              flex: 1,
              display: 'flex',
              flexDirection: 'column',
              gap: 2,
              minHeight: { xs: '32vh', md: '55vh' },
              height: { md: '70vh' },
              backgroundColor: 'rgba(148, 163, 184, 0.08)',
              borderRadius: 2,
              border: '1px solid rgba(148, 163, 184, 0.24)',
              padding: { xs: 2, md: 3 },
              boxShadow: '0 10px 30px rgba(15, 23, 42, 0.18)',
            }}
          >
            <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
              {t('documentEditor.previewTitle')}
            </Typography>
            <Box sx={{ flex: 1, overflowY: 'auto' }}>
              {content.trim() ? (
                <Box sx={markdownStyles}>
                  <ReactMarkdown remarkPlugins={[remarkGfm]}>
                    {content}
                  </ReactMarkdown>
                </Box>
              ) : (
                <Typography color="text.secondary">
                  {t('documentEditor.previewEmpty')}
                </Typography>
              )}
            </Box>
          </Box>
        </Box>
      )}
    </Box>
  );
};

export default DocumentEditor;
