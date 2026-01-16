(function () {
	'use strict';

	let activeDialogId = null; // текущий открытый диалог
	let hiddenElements = []; // элементы, скрытые через aria-hidden

	function getFocusableElements(container) {
		return Array.from(
			container.querySelectorAll(
				'button:not([disabled]), [href], input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])'
			)
		).filter(el => el.offsetParent !== null);
	}

	// скрыть все элементы кроме текущего диалога
	function hideBackground(dialogContent) {
		hiddenElements = [];
		document.body.querySelectorAll('body > *').forEach(el => {
			if (el !== dialogContent && !dialogContent.contains(el)) {
				el.setAttribute('aria-hidden', 'true');
				hiddenElements.push(el);
			}
		});
	}

	// вернуть aria-hidden всем скрытым элементам
	function restoreBackground() {
		hiddenElements.forEach(el => el.removeAttribute('aria-hidden'));
		hiddenElements = [];
	}

	function openDialog(dialogId) {
		const backdrop = document.querySelector(
			`[data-tui-dialog-backdrop][data-dialog-instance="${dialogId}"]`
		);
		const content = document.querySelector(
			`[data-tui-dialog-content][data-dialog-instance="${dialogId}"]`
		);
		if (!backdrop || !content) return;

		backdrop.removeAttribute('data-tui-dialog-hidden');
		content.removeAttribute('data-tui-dialog-hidden');

		requestAnimationFrame(() => {
			backdrop.setAttribute('data-tui-dialog-open', 'true');
			content.setAttribute('data-tui-dialog-open', 'true');
			document.body.style.overflow = 'hidden';

			// обновляем триггеры
			document
				.querySelectorAll(
					`[data-tui-dialog-trigger][data-dialog-instance="${dialogId}"]`
				)
				.forEach(trigger => {
					trigger.setAttribute('data-tui-dialog-trigger-open', 'true');
				});

			// a11y

			hideBackground(content);

			const disableAutoFocus = content.hasAttribute(
				'data-tui-dialog-disable-autofocus'
			);
			if (!disableAutoFocus) {
				setTimeout(() => {
					const focusable = getFocusableElements(content);
					if (focusable.length) focusable[0].focus();
				}, 50);
			}

			activeDialogId = dialogId;
		});
	}

	function closeDialog(dialogId) {
		const backdrop = document.querySelector(
			`[data-tui-dialog-backdrop][data-dialog-instance="${dialogId}"]`
		);
		const content = document.querySelector(
			`[data-tui-dialog-content][data-dialog-instance="${dialogId}"]`
		);
		if (!backdrop || !content) return;

		backdrop.setAttribute('data-tui-dialog-open', 'false');
		content.setAttribute('data-tui-dialog-open', 'false');

		document
			.querySelectorAll(
				`[data-tui-dialog-trigger][data-dialog-instance="${dialogId}"]`
			)
			.forEach(trigger => {
				trigger.setAttribute('data-tui-dialog-trigger-open', 'false');
			});

		setTimeout(() => {
			backdrop.setAttribute('data-tui-dialog-hidden', 'true');
			content.setAttribute('data-tui-dialog-hidden', 'true');

			const hasOpenDialogs = document.querySelector(
				'[data-tui-dialog-content][data-tui-dialog-open="true"]'
			);
			if (!hasOpenDialogs) {
				document.body.style.overflow = '';
				restoreBackground();
			}

			if (activeDialogId === dialogId) activeDialogId = null;
		}, 300);
	}

	function getDialogInstance(element) {
		const instance = element.getAttribute('data-dialog-instance');
		if (instance) return instance;
		const parentDialog = element.closest('[data-tui-dialog]');
		return parentDialog
			? parentDialog.getAttribute('data-dialog-instance')
			: null;
	}

	function isDialogOpen(dialogId) {
		const content = document.querySelector(
			`[data-tui-dialog-content][data-dialog-instance="${dialogId}"]`
		);
		return content?.getAttribute('data-tui-dialog-open') === 'true' || false;
	}

	function toggleDialog(dialogId) {
		isDialogOpen(dialogId) ? closeDialog(dialogId) : openDialog(dialogId);
	}

	// Click delegation
	document.addEventListener('click', e => {
		const trigger = e.target.closest('[data-tui-dialog-trigger]');
		if (trigger) {
			const dialogId = trigger.getAttribute('data-dialog-instance');
			if (!dialogId) return;
			toggleDialog(dialogId);
			return;
		}

		const closeBtn = e.target.closest('[data-tui-dialog-close]');
		if (closeBtn) {
			const forValue = closeBtn.getAttribute('data-tui-dialog-close');
			const dialogId = forValue || getDialogInstance(closeBtn);
			if (dialogId) closeDialog(dialogId);
			return;
		}

		const backdrop = e.target.closest('[data-tui-dialog-backdrop]');
		if (backdrop) {
			const dialogId = backdrop.getAttribute('data-dialog-instance');
			if (!dialogId) return;

			const wrapper = document.querySelector(
				`[data-tui-dialog][data-dialog-instance="${dialogId}"]`
			);
			const content = document.querySelector(
				`[data-tui-dialog-content][data-dialog-instance="${dialogId}"]`
			);

			const isDisabled =
				wrapper?.hasAttribute('data-tui-dialog-disable-click-away') ||
				content?.hasAttribute('data-tui-dialog-disable-click-away');

			if (!isDisabled) closeDialog(dialogId);
		}
	});

	// ESC
	document.addEventListener('keydown', e => {
		if (e.key === 'Escape') {
			const openDialogs = document.querySelectorAll(
				'[data-tui-dialog-content][data-tui-dialog-open="true"]'
			);
			if (openDialogs.length === 0) return;

			const content = openDialogs[openDialogs.length - 1];
			const dialogId = content.getAttribute('data-dialog-instance');
			if (!dialogId) return;

			const wrapper = document.querySelector(
				`[data-tui-dialog][data-dialog-instance="${dialogId}"]`
			);
			const isDisabled =
				wrapper?.hasAttribute('data-tui-dialog-disable-esc') ||
				content?.hasAttribute('data-tui-dialog-disable-esc');

			if (!isDisabled) closeDialog(dialogId);
		}
	});

	// Tab focus trap
	document.addEventListener('keydown', e => {
		if (e.key !== 'Tab') return;
		if (!activeDialogId) return;

		const content = document.querySelector(
			`[data-tui-dialog-content][data-dialog-instance="${activeDialogId}"]`
		);
		if (!content) return;

		const focusable = getFocusableElements(content);
		if (focusable.length === 0) return;

		const first = focusable[0];
		const last = focusable[focusable.length - 1];
		const active = document.activeElement;

		if (e.shiftKey) {
			if (active === first || !content.contains(active)) {
				e.preventDefault();
				last.focus();
			}
		} else {
			if (active === last) {
				e.preventDefault();
				first.focus();
			}
		}
	});

	document.addEventListener('DOMContentLoaded', () => {
		const openDialogs = document.querySelectorAll(
			'[data-tui-dialog-content][data-tui-dialog-open="true"]'
		);
		if (openDialogs.length > 0) document.body.style.overflow = 'hidden';
	});

	const observer = new MutationObserver(() => {
		const hasOpenDialogs = document.querySelector(
			'[data-tui-dialog-content][data-tui-dialog-open="true"]'
		);
		if (!hasOpenDialogs) document.body.style.overflow = '';
	});
	observer.observe(document.body, { childList: true, subtree: true });

	// Public API
	window.tui = window.tui || {};
	window.tui.dialog = {
		open: openDialog,
		close: closeDialog,
		toggle: toggleDialog,
		isOpen: isDialogOpen,
	};
})();
