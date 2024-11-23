// hooks/useClimatiq.ts
import { useState } from "react";
import { ClimatiqAPI } from "@/services/climatiq/api";
import type {
  SelectorModel,
  ParametersModel,
  EstimationModel,
  EmissionFactorResponse,
  EstimationResponse,
} from "@/services/climatiq/types/models";

export function useClimatiq() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  // Selector Hook
  async function searchEmissionFactors(params: SelectorModel) {
    try {
      setLoading(true);
      setError(null);
      return await ClimatiqAPI.searchEmissionFactors(params);
    } catch (err) {
      setError(err instanceof Error ? err : new Error("Search failed"));
      throw err;
    } finally {
      setLoading(false);
    }
  }

  // Parameters Hook
  async function getEmissionFactors(params: ParametersModel) {
    try {
      setLoading(true);
      setError(null);
      return await ClimatiqAPI.getEmissionFactors(params);
    } catch (err) {
      setError(
        err instanceof Error ? err : new Error("Failed to get emission factors")
      );
      throw err;
    } finally {
      setLoading(false);
    }
  }

  // Estimation Hook
  async function calculateEmissions(params: EstimationModel) {
    try {
      setLoading(true);
      setError(null);
      return await ClimatiqAPI.calculateEmissions(params);
    } catch (err) {
      setError(err instanceof Error ? err : new Error("Calculation failed"));
      throw err;
    } finally {
      setLoading(false);
    }
  }

  return {
    loading,
    error,
    searchEmissionFactors,
    getEmissionFactors,
    calculateEmissions,
  };
}