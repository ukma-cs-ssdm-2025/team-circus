export type CursorLocation = {
	line: number;
	column: number;
};

export type RemoteUserPresence = {
	id: string;
	name: string;
	role?: string;
	status?: "online" | "away" | "offline";
	color?: string;
	cursorPosition?: number;
	cursorLocation?: CursorLocation;
};

export type EditorProps = {
	value?: string;
	defaultValue?: string;
	onChange?: (nextValue: string) => void;
	onCursorChange?: (position: number) => void;
	remoteUsers?: RemoteUserPresence[];
	isConnected?: boolean;
	showPresence?: boolean;
	className?: string;
	ariaLabel?: string;
	colorScheme?: "light" | "dark";
};
