import { Card, CardContent, Chip, Stack, Typography } from '@mui/material';
import { formatDate, truncateText } from '../../utils';
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
  return (
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
    </Card>
  );
};

export default DocumentCard;
