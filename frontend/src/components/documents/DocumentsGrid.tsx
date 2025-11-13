import Grid from "@mui/material/Grid";
import DocumentCard from "./DocumentCard";
import type { DocumentItem } from "../../types/entities";

interface DocumentPermissions {
  canEdit: boolean;
  canDelete: boolean;
}

interface DocumentsGridProps {
  documents: DocumentItem[];
  groupNameByUUID: Record<string, string>;
  createdAtLabel: string;
  noContentLabel: string;
  groupUnknownLabel: string;
  editLabel: string;
  deleteLabel: string;
  permissionsByDocument?: Record<string, DocumentPermissions>;
  onDocumentDelete?: (document: DocumentItem) => void;
}

const DocumentsGrid = ({
  documents,
  groupNameByUUID,
  createdAtLabel,
  noContentLabel,
  groupUnknownLabel,
  editLabel,
  deleteLabel,
  permissionsByDocument,
  onDocumentDelete,
}: DocumentsGridProps) => {
  return (
    <Grid container spacing={3}>
      {documents.map((document) => {
        const documentPermissions = permissionsByDocument?.[document.uuid] ?? {
          canEdit: true,
          canDelete: true,
        };

        return (
          <Grid key={document.uuid} size={{ xs: 12, md: 6 }}>
            <DocumentCard
              document={document}
              groupName={groupNameByUUID[document.group_uuid]}
              createdAtLabel={createdAtLabel}
              noContentLabel={noContentLabel}
              groupUnknownLabel={groupUnknownLabel}
              editLabel={editLabel}
              deleteLabel={deleteLabel}
              canEdit={documentPermissions.canEdit}
              canDelete={documentPermissions.canDelete}
              onDelete={onDocumentDelete}
            />
          </Grid>
        );
      })}
    </Grid>
  );
};

export default DocumentsGrid;
