import Alpine from 'alpinejs';
import focus from '@alpinejs/focus';

import 'htmx.org';

import { A } from '../shared/keys';
import { Modal, MODAL_COMPONENT_DATA } from '@/components/modal/modal';

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

function Hello() {
	console.log(`Variable from other file: ${A}`);
}

Alpine.plugin(focus);

document.addEventListener('alpine:init', () => {
	Alpine.data<ToastData, ToastDataInitialArgs>('toast', () => ({
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

	Alpine.data(MODAL_COMPONENT_DATA, (initialOpen: boolean) => {
		return Modal(initialOpen);
	});
});

Alpine.start();

Hello();
