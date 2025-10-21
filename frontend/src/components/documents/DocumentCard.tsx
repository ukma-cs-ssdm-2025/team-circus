import { Card, CardActionArea, CardContent, Chip, Stack, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import { formatDate, truncateText } from '../../utils';
import { ROUTES } from '../../constants';
import type { DocumentItem } from '../../types/entities';

interface DocumentCardProps {
  document: DocumentItem;
  groupName: string;
  createdAtLabel: string;
  noContentLabel: string;
  groupUnknownLabel: string;
}

const DocumentCard = ({
  document,
  groupName,
  createdAtLabel,
  noContentLabel,
  groupUnknownLabel
}: DocumentCardProps) => {
  const documentPath = `${ROUTES.DOCUMENTS}/${document.uuid}`;

  return (
    <Card
      sx={{
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        borderRadius: 3,
        boxShadow: '0 10px 30px rgba(0, 0, 0, 0.08)',
        transition: 'transform 0.2s ease, box-shadow 0.2s ease',
        '&:hover': {
          transform: 'translateY(-4px)',
          boxShadow: '0 16px 32px rgba(0, 0, 0, 0.12)',
        },
      }}
    >
      <CardActionArea
        component={RouterLink}
        to={documentPath}
        sx={{ height: '100%', alignItems: 'stretch', display: 'flex' }}
      >
        <CardContent sx={{ flex: 1 }}>
          <Stack spacing={2}>
            <Stack direction="row" spacing={1} alignItems="center">
              <Typography variant="h6" sx={{ fontWeight: 600 }}>
                {document.name}
              </Typography>
              <Chip
                label={groupName || groupUnknownLabel}
                size="small"
                color="primary"
                sx={{ fontWeight: 600 }}
              />
            </Stack>

            <Typography variant="body2" color="text.secondary">
              {truncateText(document.content || noContentLabel, 180)}
            </Typography>

            <Typography variant="body2" color="text.secondary">
              {`${createdAtLabel}: ${formatDate(document.created_at)}`}
            </Typography>
          </Stack>
        </CardContent>
      </CardActionArea>
    </Card>
  );
};

export default DocumentCard;
