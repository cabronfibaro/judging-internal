export type Option = {
	label: string;
	value: number;
};

export type Parameter = {
	name: string;
	type: string;
};

export type Tag = {
	id: string;
	label: string;
	type: 'bonus' | 'deduction' | 'neutral';
	value_type: 'dynamic' | 'enum' | 'switch';
	options?: Option[];
	default_value?: number;
	parameter?: Parameter;
	max?: number;
	multiplier?: number;
};
