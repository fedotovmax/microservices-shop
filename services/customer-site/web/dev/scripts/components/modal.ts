import Alpine from 'alpinejs';

export const MODAL_COMPONENT_DATA = 'modal';

export function initModalComponent() {
	Alpine.data(MODAL_COMPONENT_DATA, (initialOpen: boolean) => ({
		isOpen: initialOpen,
	}));
}
