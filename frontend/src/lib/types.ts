export type AvatarMe = {
	name: string;
	email: string;
	imageUri?: string | null;
	isSystemAdmin: boolean;
	isAnalyst: boolean;
	canPairing: boolean;
};

export type NavItem = {
	href: string;
	label: string;
	icon: string;
	active: boolean;
};
