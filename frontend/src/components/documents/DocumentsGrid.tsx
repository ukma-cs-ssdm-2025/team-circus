import { Grid } from '@mui/material';
import DocumentCard from './DocumentCard';
import type { DocumentItem } from '../../types/entities';

interface DocumentsGridProps {
  documents: DocumentItem[];
  groupNameByUUID: Record<string, string>;
  createdAtLabel: string;
  noContentLabel: string;
  groupUnknownLabel: string;
}

const DocumentsGrid = ({ 
  documents, 
  groupNameByUUID, 
  createdAtLabel, 
  noContentLabel, 
  groupUnknownLabel 
}: DocumentsGridProps) => {
  return (
    <Grid container spacing={3}>
      {documents.map(document => (
        <Grid item xs={12} md={6} key={document.uuid}>
          <DocumentCard
            document={document}
            groupName={groupNameByUUID[document.group_uuid]}
            createdAtLabel={createdAtLabel}
            noContentLabel={noContentLabel}
            groupUnknownLabel={groupUnknownLabel}
          />
        </Grid>
      ))}
    </Grid>
  );
};

export default DocumentsGrid;
