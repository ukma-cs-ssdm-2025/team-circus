import type { MouseEvent } from "react";
import {
  Card,
  CardActionArea,
  CardActions,
  CardContent,
  Chip,
  Stack,
  Typography,
  Button,
} from "@mui/material";
import { Link as RouterLink } from "react-router-dom";
import { formatDate, truncateText } from "../../utils";
import { ROUTES } from "../../constants";
import type { DocumentItem } from "../../types/entities";

interface DocumentCardProps {
  document: DocumentItem;
  groupName: string;
  createdAtLabel: string;
  noContentLabel: string;
  groupUnknownLabel: string;
  editLabel: string;
  deleteLabel: string;
  canEdit?: boolean;
  canDelete?: boolean;
  onDelete?: (document: DocumentItem) => void;
}

const DocumentCard = ({
  document,
  groupName,
  createdAtLabel,
  noContentLabel,
  groupUnknownLabel,
  editLabel,
  deleteLabel,
  canEdit = true,
  canDelete = true,
  onDelete,
}: DocumentCardProps) => {
  const documentPath = `${ROUTES.DOCUMENTS}/${document.uuid}`;

  const handleDeleteClick = (event: MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();
    event.stopPropagation();
    onDelete?.(document);
  };

  return (
    <Card
      sx={{
        height: "100%",
        display: "flex",
        flexDirection: "column",
        borderRadius: 3,
        boxShadow: "0 10px 30px rgba(0, 0, 0, 0.08)",
        transition: "transform 0.2s ease, box-shadow 0.2s ease",
        "&:hover": {
          transform: "translateY(-4px)",
          boxShadow: "0 16px 32px rgba(0, 0, 0, 0.12)",
        },
      }}
    >
      <CardActionArea
        component={RouterLink}
        to={documentPath}
        sx={{ height: "100%", alignItems: "stretch", display: "flex" }}
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
      {(canEdit || canDelete) && (
        <CardActions sx={{ justifyContent: "flex-end", pt: 0, pb: 2, px: 2 }}>
          {canEdit && (
            <Button
              size="small"
              component={RouterLink}
              to={documentPath}
              onClick={(event) => event.stopPropagation()}
            >
              {editLabel}
            </Button>
          )}
          {canDelete && (
            <Button
              size="small"
              color="error"
              onClick={handleDeleteClick}
              disabled={!onDelete}
            >
              {deleteLabel}
            </Button>
          )}
        </CardActions>
      )}
    </Card>
  );
};

export default DocumentCard;
