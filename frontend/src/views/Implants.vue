<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { ArrowLeft, Loader2 } from 'lucide-vue-next'
import OperationCard from '@/components/OperationCard.vue'
import type { Operations } from '@/types/implantTypes'
import { SavePatient, GetAllPatients } from '@wailsjs/go/excel/ExcelService'
import { excel } from '@wailsjs/go/models'
import { toast } from 'vue-sonner'

const router = useRouter()

// Form data
const patientName = ref<string>('')
const implantCount = ref<number | string>('')

// Operations data
const operations = ref<Operations>({
  removal: {
    dates: [],
    numbers: [],
    title: 'Removal',
    comment: '',
  },
  sinusLift: {
    dates: [],
    numbers: [],
    title: 'Sinus Lift',
    comment: '',
  },
  boneGrafting: {
    dates: [],
    numbers: [],
    title: 'Bone Grafting',
    comment: '',
  },
  installationFormation: {
    dates: [],
    numbers: [],
    title: 'Installation Formation',
    comment: '',
  },
  reinstallationImplant: {
    dates: [],
    numbers: [],
    title: 'Reinstallation Implant',
    comment: '',
  },
  permanentProsthetics: {
    dates: [],
    numbers: [],
    title: 'Permanent Prosthetics',
    comment: '',
  },
  temporaryProsthetics: {
    dates: [],
    numbers: [],
    title: 'Temporary Prosthetics',
    comment: '',
  },
})

// Update operation handlers
const updateOperationDates = (operationKey: keyof Operations, dates: string[]) => {
  operations.value[operationKey].dates = dates as any
}

const updateOperationNumbers = (operationKey: keyof Operations, numbers: string[]) => {
  operations.value[operationKey].numbers = numbers
}

const updateOperationComment = (operationKey: keyof Operations, comment: string) => {
  operations.value[operationKey].comment = comment
}

// Patient list and loading state
const patients = ref<excel.Patient[]>([])
const isLoading = ref(false)
const isSaving = ref(false)
const activeTab = ref('add-patient')

// Load patients when switching to patient list tab
const loadPatients = async () => {
  isLoading.value = true
  try {
    patients.value = await GetAllPatients()
  } catch (error) {
    console.error('Failed to load patients:', error)
    toast.error('Failed to load patients', {
      description: String(error)
    })
  } finally {
    isLoading.value = false
  }
}

// Handle tab change
const handleTabChange = (value: string | number) => {
  const tabValue = String(value)
  activeTab.value = tabValue
  if (tabValue === 'patient-list') {
    loadPatients()
  }
}

// Form submission
const handleSubmit = async () => {
  // Validation
  if (!patientName.value.trim()) {
    toast.error('Validation Error', {
      description: 'Please enter patient name'
    })
    return
  }

  if (!implantCount.value || Number(implantCount.value) <= 0) {
    toast.error('Validation Error', {
      description: 'Please enter a valid implant count'
    })
    return
  }

  isSaving.value = true
  try {
    // Prepare patient data for backend
    const patientData = {
      id: '',
      fio: patientName.value.trim(),
      implantNumber: Number(implantCount.value),
      operations: {
        removal: {
          dates: operations.value.removal.dates,
          numbers: operations.value.removal.numbers,
          comment: operations.value.removal.comment || '',
        },
        sinusLift: {
          dates: operations.value.sinusLift.dates,
          numbers: operations.value.sinusLift.numbers,
          comment: operations.value.sinusLift.comment || '',
        },
        boneGrafting: {
          dates: operations.value.boneGrafting.dates,
          numbers: operations.value.boneGrafting.numbers,
          comment: operations.value.boneGrafting.comment || '',
        },
        installationFormation: {
          dates: operations.value.installationFormation.dates,
          numbers: operations.value.installationFormation.numbers,
          comment: operations.value.installationFormation.comment || '',
        },
        reinstallationImplant: {
          dates: operations.value.reinstallationImplant.dates,
          numbers: operations.value.reinstallationImplant.numbers,
          comment: operations.value.reinstallationImplant.comment || '',
        },
        permanentProsthetics: {
          dates: operations.value.permanentProsthetics.dates,
          numbers: operations.value.permanentProsthetics.numbers,
          comment: operations.value.permanentProsthetics.comment || '',
        },
        temporaryProsthetics: {
          dates: operations.value.temporaryProsthetics.dates,
          numbers: operations.value.temporaryProsthetics.numbers,
          comment: operations.value.temporaryProsthetics.comment || '',
        },
      },
      controlHalfYear: '',
      controlYear: '',
      occupationalHygiene: '',
    }

    await SavePatient(excel.Patient.createFrom(patientData))

    toast.success('Success', {
      description: 'Patient saved successfully!'
    })

    // Reset form
    patientName.value = ''
    implantCount.value = ''

    // Reset operations
    Object.keys(operations.value).forEach((key) => {
      const opKey = key as keyof Operations
      operations.value[opKey].dates = []
      operations.value[opKey].numbers = []
      operations.value[opKey].comment = ''
    })
  } catch (error) {
    console.error('Failed to save patient:', error)
    toast.error('Failed to save patient', {
      description: String(error)
    })
  } finally {
    isSaving.value = false
  }
}
</script>

<template>
  <div class="implants-page h-screen flex flex-col bg-gradient-to-br from-gray-50 via-white to-gray-50">
    <!-- Header -->
    <header class="bg-white border-b border-gray-200 flex-shrink-0">
      <div class="px-4 py-3 flex items-center gap-3">
        <Button
          variant="outline"
          size="icon"
          @click="router.push('/')"
          class="flex-shrink-0 border-gray-300 hover:bg-blue-50 hover:border-blue-500"
        >
          <ArrowLeft class="h-5 w-5 text-gray-700" />
        </Button>
        <div class="flex-1">
          <h1 class="text-xl font-bold text-gray-900">Implants Management</h1>
          <p class="text-xs text-gray-600 mt-0.5">Manage dental implants and track procedures</p>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex-1 overflow-auto px-4 py-4">
      <Tabs default-value="add-patient" class="w-full" @update:model-value="handleTabChange">
        <TabsList class="grid w-full grid-cols-2 max-w-md">
          <TabsTrigger value="add-patient">Add Patient</TabsTrigger>
          <TabsTrigger value="patient-list">Patient List</TabsTrigger>
        </TabsList>

        <!-- Add Patient Tab -->
        <TabsContent value="add-patient" class="mt-4">
          <Card>
            <CardHeader>
              <CardTitle>Patient Information</CardTitle>
            </CardHeader>
            <CardContent class="space-y-4">
              <!-- Patient Details -->
              <div class="grid grid-cols-2 gap-3">
                <div class="space-y-2">
                  <Label for="patient-name">Full Name</Label>
                  <Input
                    id="patient-name"
                    v-model="patientName"
                    placeholder="Enter patient full name"
                  />
                </div>
                <div class="space-y-2">
                  <Label for="implant-count">Implant Count</Label>
                  <Input
                    id="implant-count"
                    v-model="implantCount"
                    type="number"
                    placeholder="Enter number of implants"
                  />
                </div>
              </div>

              <!-- Operations Grid -->
              <div class="space-y-3">
                <h3 class="text-base font-semibold text-gray-900">Operations</h3>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                  <OperationCard
                    :operation-data="operations.removal"
                    @update:dates="(dates: string[]) => updateOperationDates('removal', dates)"
                    @update:numbers="(numbers: string[]) => updateOperationNumbers('removal', numbers)"
                    @update:comment="(comment: string) => updateOperationComment('removal', comment)"
                  />
                  <OperationCard
                    :operation-data="operations.sinusLift"
                    @update:dates="(dates: string[]) => updateOperationDates('sinusLift', dates)"
                    @update:numbers="(numbers: string[]) => updateOperationNumbers('sinusLift', numbers)"
                    @update:comment="(comment: string) => updateOperationComment('sinusLift', comment)"
                  />
                  <OperationCard
                    :operation-data="operations.boneGrafting"
                    @update:dates="(dates: string[]) => updateOperationDates('boneGrafting', dates)"
                    @update:numbers="(numbers: string[]) => updateOperationNumbers('boneGrafting', numbers)"
                    @update:comment="(comment: string) => updateOperationComment('boneGrafting', comment)"
                  />
                  <OperationCard
                    :operation-data="operations.installationFormation"
                    @update:dates="(dates: string[]) => updateOperationDates('installationFormation', dates)"
                    @update:numbers="
                      (numbers: string[]) => updateOperationNumbers('installationFormation', numbers)
                    "
                    @update:comment="
                      (comment: string) => updateOperationComment('installationFormation', comment)
                    "
                  />
                  <OperationCard
                    :operation-data="operations.reinstallationImplant"
                    @update:dates="(dates: string[]) => updateOperationDates('reinstallationImplant', dates)"
                    @update:numbers="
                      (numbers: string[]) => updateOperationNumbers('reinstallationImplant', numbers)
                    "
                    @update:comment="
                      (comment: string) => updateOperationComment('reinstallationImplant', comment)
                    "
                  />
                  <OperationCard
                    :operation-data="operations.permanentProsthetics"
                    @update:dates="(dates: string[]) => updateOperationDates('permanentProsthetics', dates)"
                    @update:numbers="
                      (numbers: string[]) => updateOperationNumbers('permanentProsthetics', numbers)
                    "
                    @update:comment="
                      (comment: string) => updateOperationComment('permanentProsthetics', comment)
                    "
                  />
                  <OperationCard
                    :operation-data="operations.temporaryProsthetics"
                    @update:dates="(dates: string[]) => updateOperationDates('temporaryProsthetics', dates)"
                    @update:numbers="
                      (numbers: string[]) => updateOperationNumbers('temporaryProsthetics', numbers)
                    "
                    @update:comment="
                      (comment: string) => updateOperationComment('temporaryProsthetics', comment)
                    "
                  />
                </div>
              </div>

              <!-- Submit Button -->
              <div class="flex justify-end pt-3">
                <Button @click="handleSubmit" size="lg" :disabled="isSaving">
                  <Loader2 v-if="isSaving" class="mr-2 h-4 w-4 animate-spin" />
                  {{ isSaving ? 'Saving...' : 'Save Patient' }}
                </Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <!-- Patient List Tab -->
        <TabsContent value="patient-list" class="mt-4">
          <Card>
            <CardHeader>
              <CardTitle>Patient List</CardTitle>
            </CardHeader>
            <CardContent>
              <div class="rounded-md border overflow-auto max-h-[450px]">
                <Table>
                  <TableHeader class="sticky top-0 bg-muted z-10">
                    <TableRow>
                      <TableHead class="w-[80px]">ID</TableHead>
                      <TableHead class="w-[200px]">Full Name</TableHead>
                      <TableHead class="w-[100px] text-center">Implants</TableHead>
                      <TableHead class="w-[150px]">Removal</TableHead>
                      <TableHead class="w-[150px]">Sinus Lift</TableHead>
                      <TableHead class="w-[150px]">Bone Grafting</TableHead>
                      <TableHead class="w-[150px]">Formation</TableHead>
                      <TableHead class="w-[150px]">Reinstallation</TableHead>
                      <TableHead class="w-[150px]">Permanent Pros.</TableHead>
                      <TableHead class="w-[150px]">Temporary Pros.</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    <TableRow v-if="isLoading">
                      <TableCell colspan="10" class="h-24 text-center text-muted-foreground">
                        <Loader2 class="h-6 w-6 animate-spin mx-auto" />
                        <p class="mt-2">Loading patients...</p>
                      </TableCell>
                    </TableRow>
                    <TableRow v-else-if="patients.length === 0">
                      <TableCell colspan="10" class="h-24 text-center text-muted-foreground">
                        No patients found. Add a patient to get started.
                      </TableCell>
                    </TableRow>
                    <TableRow v-else v-for="patient in patients" :key="patient.id">
                      <TableCell class="font-medium">{{ patient.id }}</TableCell>
                      <TableCell>{{ patient.fio }}</TableCell>
                      <TableCell class="text-center">{{ patient.implantNumber }}</TableCell>
                      <TableCell>{{ patient.operations.removal?.dates.join(', ') || '-' }}</TableCell>
                      <TableCell>{{ patient.operations.sinusLift?.dates.join(', ') || '-' }}</TableCell>
                      <TableCell>{{ patient.operations.boneGrafting?.dates.join(', ') || '-' }}</TableCell>
                      <TableCell>{{ patient.operations.installationFormation?.dates.join(', ') || '-' }}</TableCell>
                      <TableCell>{{ patient.operations.reinstallationImplant?.dates.join(', ') || '-' }}</TableCell>
                      <TableCell>{{ patient.operations.permanentProsthetics?.dates.join(', ') || '-' }}</TableCell>
                      <TableCell>{{ patient.operations.temporaryProsthetics?.dates.join(', ') || '-' }}</TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </main>
  </div>
</template>

<style scoped>
/* Ensure smooth scrolling for the table */
.overflow-auto {
  scrollbar-width: thin;
  scrollbar-color: hsl(var(--muted-foreground)) hsl(var(--muted));
}

.overflow-auto::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.overflow-auto::-webkit-scrollbar-track {
  background: hsl(var(--muted));
  border-radius: 4px;
}

.overflow-auto::-webkit-scrollbar-thumb {
  background: hsl(var(--muted-foreground));
  border-radius: 4px;
}

.overflow-auto::-webkit-scrollbar-thumb:hover {
  background: hsl(var(--foreground));
}
</style>
