import Alpine from 'alpinejs';
import focus from '@alpinejs/focus';

import 'htmx.org';

import { initToastComponent } from '@/components/toast';
import { initModalComponent } from '@/components/modal';
import { PAGE_PROPS } from '@/shared/selectors';
import { htmxRequest } from '@/shared/htmx';

document.addEventListener('DOMContentLoaded', () => {
	const pageData = document.getElementById(PAGE_PROPS);

	if (pageData) {
		try {
			const data = JSON.parse(pageData.textContent);
			console.log(data);
		} catch (error) {
			console.error(error);
		}
	}
});

Alpine.plugin(focus);

document.addEventListener('alpine:init', () => {
	initToastComponent();
	initModalComponent();

	Alpine.data('htmxnotify', () => ({
		isLoading: false,

		async notify() {
			try {
				this.isLoading = true;
				await htmxRequest('/notify', 'get');
			} catch (error) {
				console.error(error);
			} finally {
				this.isLoading = false;
			}
		},
	}));
});

Alpine.start();
