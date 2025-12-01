import Grid from "@mui/material/Grid";
import type { DocumentItem } from "../../types/entities";
import DocumentCard from "./DocumentCard";

interface DocumentsGridProps {
	documents: DocumentItem[];
	groupNameByUUID: Record<string, string>;
	createdAtLabel: string;
	noContentLabel: string;
	groupUnknownLabel: string;
	deleteLabel?: string;
	onDeleteDocument?: (document: DocumentItem) => void;
}

const DocumentsGrid = ({
	documents,
	groupNameByUUID,
	createdAtLabel,
	noContentLabel,
	groupUnknownLabel,
	deleteLabel,
	onDeleteDocument,
}: DocumentsGridProps) => {
	return (
		<Grid container spacing={3}>
			{documents.map((document) => (
				<Grid key={document.uuid} size={{ xs: 12, md: 6 }}>
					<DocumentCard
						document={document}
						groupName={groupNameByUUID[document.group_uuid]}
						createdAtLabel={createdAtLabel}
						noContentLabel={noContentLabel}
						groupUnknownLabel={groupUnknownLabel}
						deleteLabel={deleteLabel}
						onDelete={onDeleteDocument}
					/>
				</Grid>
			))}
		</Grid>
	);
};

export default DocumentsGrid;
