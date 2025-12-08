import {
	Button,
	Dialog,
	DialogActions,
	DialogContent,
	DialogContentText,
	DialogTitle,
	Stack,
} from "@mui/material";
import type { ReactNode } from "react";

interface ConfirmDialogProps {
	open: boolean;
	title: string;
	description: ReactNode;
	confirmLabel: string;
	cancelLabel: string;
	onConfirm: () => void | Promise<void>;
	onClose: () => void;
	confirming?: boolean;
}

export const ConfirmDialog = ({
	open,
	title,
	description,
	confirmLabel,
	cancelLabel,
	onConfirm,
	onClose,
	confirming = false,
}: ConfirmDialogProps) => {
	return (
		<Dialog
			open={open}
			onClose={confirming ? undefined : onClose}
			fullWidth
			maxWidth="xs"
		>
			<DialogTitle>{title}</DialogTitle>
			<DialogContent>
				{typeof description === "string" ? (
					<DialogContentText>{description}</DialogContentText>
				) : (
					<Stack spacing={1}>{description}</Stack>
				)}
			</DialogContent>
			<DialogActions>
				<Button onClick={onClose} disabled={confirming} color="inherit">
					{cancelLabel}
				</Button>
				<Button
					onClick={onConfirm}
					color="error"
					disabled={confirming}
					variant="contained"
				>
					{confirmLabel}
				</Button>
			</DialogActions>
		</Dialog>
	);
};
