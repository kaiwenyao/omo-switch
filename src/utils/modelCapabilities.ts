export type AgentVariant = 'none' | 'low' | 'medium' | 'high' | 'xhigh' | 'max';
export type ActiveVariant = Exclude<AgentVariant, 'none'>;
export type VariantProfileKind = 'legacy' | 'openai-reasoning';

export interface VariantProfile {
  kind: VariantProfileKind;
  allowedVariants: ActiveVariant[];
  supportsUnset: boolean;
}

const LEGACY_PROFILE: VariantProfile = {
  kind: 'legacy',
  allowedVariants: ['max', 'high', 'medium', 'low'],
  supportsUnset: true,
};

const OPENAI_CLONE_PROVIDERS = new Set([
  'aicodewith',
]);

const OPENAI_REASONING_RULES: Array<{
  test: (modelId: string) => boolean;
  allowedVariants: ActiveVariant[];
}> = [
  {
    test: (modelId) => /^gpt-5(?:\.0)?$/.test(modelId),
    allowedVariants: ['high', 'medium', 'low'],
  },
  {
    test: (modelId) => /^gpt-5\.1(?:-mini)?$/.test(modelId),
    allowedVariants: ['high', 'medium', 'low'],
  },
  {
    test: (modelId) => /^gpt-5-pro$/.test(modelId),
    allowedVariants: ['high'],
  },
  {
    test: (modelId) => /^gpt-5\.(?:2|[3-9]|\d{2,})(?:-codex)?$/.test(modelId),
    allowedVariants: ['xhigh', 'high', 'medium', 'low'],
  },
  {
    test: (modelId) => /^gpt-5\.(?:2|[3-9]|\d{2,})-pro$/.test(modelId),
    allowedVariants: ['xhigh', 'high', 'medium'],
  },
];

const VARIANT_DISPLAY_ORDER: AgentVariant[] = ['none', 'low', 'medium', 'high', 'xhigh', 'max'];

function splitModelPath(modelPath: string): { providerId: string; modelId: string } {
  const [providerId = '', ...rest] = modelPath.split('/');
  return {
    providerId,
    modelId: rest.join('/'),
  };
}

function isOpenAIReasoningProvider(providerId: string): boolean {
  return providerId === 'openai' || OPENAI_CLONE_PROVIDERS.has(providerId);
}

export function getVariantProfile(modelPath: string): VariantProfile {
  const { providerId, modelId } = splitModelPath(modelPath);

  if (!providerId || !modelId || !isOpenAIReasoningProvider(providerId)) {
    return LEGACY_PROFILE;
  }

  const matchedRule = OPENAI_REASONING_RULES.find((rule) => rule.test(modelId));
  if (!matchedRule) {
    return LEGACY_PROFILE;
  }

  return {
    kind: 'openai-reasoning',
    allowedVariants: matchedRule.allowedVariants,
    supportsUnset: true,
  };
}

export function getVariantOptions(modelPath: string): AgentVariant[] {
  const profile = getVariantProfile(modelPath);
  const options = profile.supportsUnset
    ? ['none', ...profile.allowedVariants]
    : [...profile.allowedVariants];

  return VARIANT_DISPLAY_ORDER.filter((variant) => options.includes(variant));
}

export function normalizeVariantForModel(modelPath: string, variant?: string): AgentVariant {
  if (!variant || variant === 'none') {
    return 'none';
  }

  const profile = getVariantProfile(modelPath);
  const value = variant as AgentVariant;
  if (profile.kind === 'openai-reasoning' && value === 'max') {
    return profile.allowedVariants.includes('xhigh') ? 'xhigh' : profile.allowedVariants[0];
  }

  if (profile.allowedVariants.includes(value as ActiveVariant)) {
    return value;
  }

  return profile.supportsUnset ? 'none' : profile.allowedVariants[0];
}

export function getVariantDisplayValue(modelPath: string, variant?: string): AgentVariant {
  return normalizeVariantForModel(modelPath, variant);
}

export function isVariantSupported(modelPath: string, variant?: string): boolean {
  return normalizeVariantForModel(modelPath, variant) === (variant ?? 'none');
}
