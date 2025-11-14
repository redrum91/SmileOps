export interface Operation {
  dates: string[]
  numbers: string[]
  title?: string
  comment?: string
}

export interface Operations {
  removal: Operation
  sinusLift: Operation
  boneGrafting: Operation
  installationFormation: Operation
  reinstallationImplant: Operation
  permanentProsthetics: Operation
  temporaryProsthetics: Operation
}

export interface Patient {
  id: string
  fio: string
  implantNumber: number
  controlHalfYear: string
  controlYear: string
  occupationalHygiene: string
  operations: Operations
}
