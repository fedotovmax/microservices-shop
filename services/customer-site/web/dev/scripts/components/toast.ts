import Alpine from 'alpinejs';

type NotificationInput = {
	variant: string;
	sender?: string;
	title?: string;
	message?: string;
};

type Notification = {
	id: number;
	variant: string;
	sender?: string | null;
	title?: string | null;
	message?: string | null;
};

type ToastData = {
	notifications: Notification[];
	displayDuration: number;
	addNotification: (input: NotificationInput) => void;
	removeNotification: (id: number) => void;
};

type ToastDataInitialArgs = any[];

export const TOAST_COMPONENT_DATA = 'toast';

export function initToastComponent() {
	Alpine.data<ToastData, ToastDataInitialArgs>(TOAST_COMPONENT_DATA, () => ({
		notifications: [],
		displayDuration: 4000,
		addNotification(input: NotificationInput) {
			const {
				variant = 'info',
				sender = null,
				title = null,
				message = null,
			} = input;
			const id = Date.now();
			const notification = { id, variant, sender, title, message };
			if (this.notifications.length >= 20) {
				this.notifications.splice(0, this.notifications.length - 19);
			}
			this.notifications.push(notification);
		},
		removeNotification(id) {
			setTimeout(() => {
				this.notifications = this.notifications.filter(
					notification => notification.id !== id
				);
			}, 400);
		},
	}));
}
