// app/components/UnitTypesDisplay/index.tsx
/* eslint-disable */
"use client";

import React, { useEffect, useState } from "react";
import { getUnitTypes, UnitType } from "@/app/services/api";

interface UnitTypesDisplayProps {
  onUnitTypeSelect?: (unitType: string, unit: string) => void;
}

const UnitTypesDisplay: React.FC<UnitTypesDisplayProps> = ({
  onUnitTypeSelect,
}) => {
  const [unitTypes, setUnitTypes] = useState<UnitType[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedUnitType, setSelectedUnitType] = useState<string>("");
  const [selectedUnit, setSelectedUnit] = useState<string>("");

  useEffect(() => {
    const fetchUnitTypes = async () => {
      try {
        setLoading(true);
        const data = await getUnitTypes();
        // Ensure data is an array and has the correct structure
        if (Array.isArray(data)) {
          setUnitTypes(data);
        } else {
          throw new Error("Invalid data format received");
        }
        setError(null);
      } catch (err) {
        setError(
          err instanceof Error ? err.message : "Failed to fetch unit types"
        );
        setUnitTypes([]); // Reset to empty array on error
      } finally {
        setLoading(false);
      }
    };

    fetchUnitTypes();
  }, []);

  const handleUnitTypeChange = (
    event: React.ChangeEvent<HTMLSelectElement>
  ) => {
    const newUnitType = event.target.value;
    setSelectedUnitType(newUnitType);
    setSelectedUnit("");
  };

  const handleUnitChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    const newUnit = event.target.value;
    setSelectedUnit(newUnit);
    if (onUnitTypeSelect && selectedUnitType) {
      onUnitTypeSelect(selectedUnitType, newUnit);
    }
  };

  if (loading) {
    return <div>Loading unit types...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  const selectedTypeData = unitTypes.find(
    (type) => type.unit_type === selectedUnitType
  );
  const availableUnits = selectedTypeData?.units || [];

  return (
    <div className="unit-types-display">
      <div className="unit-type-selector">
        <label className="form-label" htmlFor="unitType">
          Unit Type:
        </label>
        <select
          id="unitType"
          className="form-select"
          value={selectedUnitType}
          onChange={handleUnitTypeChange}
        >
          <option value="">Select a unit type</option>
          {unitTypes.map((type) => (
            <option key={type.unit_type} value={type.unit_type}>
              {type.unit_type}
            </option>
          ))}
        </select>
      </div>

      {selectedUnitType && (
        <div className="unit-selector">
          <label className="form-label" htmlFor="unit">
            Unit:
          </label>
          <select
            id="unit"
            className="form-select"
            value={selectedUnit}
            onChange={handleUnitChange}
          >
            <option value="">Select a unit</option>
            {Array.isArray(availableUnits) &&
              availableUnits.map((unit) => (
                <option key={unit} value={unit}>
                  {unit}
                </option>
              ))}
          </select>
        </div>
      )}
    </div>
  );
};

export default UnitTypesDisplay;