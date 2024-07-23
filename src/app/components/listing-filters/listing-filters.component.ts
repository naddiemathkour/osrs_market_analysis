import { Component, Input, EventEmitter, Output } from '@angular/core';
import {
  MatButtonToggleChange,
  MatButtonToggleModule,
} from '@angular/material/button-toggle';
import { MatSelectModule } from '@angular/material/select';
import {
  MatCheckboxChange,
  MatCheckboxModule,
} from '@angular/material/checkbox';
import { MatInputModule } from '@angular/material/input';
import { IFilters } from '../../interfaces/filters.interface';
import { MatProgressBarModule } from '@angular/material/progress-bar';

@Component({
  selector: 'app-listing-filters',
  standalone: true,
  imports: [
    MatButtonToggleModule,
    MatSelectModule,
    MatCheckboxModule,
    MatInputModule,
    MatProgressBarModule,
  ],
  templateUrl: './listing-filters.component.html',
  styleUrl: './listing-filters.component.scss',
})
export class ListingFiltersComponent {
  @Input() filterValues!: IFilters;
  @Input() loading!: boolean;
  @Output() updateFilter: EventEmitter<IFilters> = new EventEmitter<IFilters>();
  newFilters!: IFilters;
  resetFilters!: IFilters;
  selectedFilter: string | null = null;

  filterNames = {
    buyLimit: 'Buy Limit',
    margin: 'Profit Margin',
    alchprof: 'Alch Profit',
    highVolume: 'High Volume',
    lowVolume: 'Low Volume',
  };

  toggleView(event: MatButtonToggleChange) {
    this.filterValues.dataType = event.value;
    this.clearFilters();
  }

  toggleMembers(event: MatCheckboxChange) {
    this.filterValues.members = !event.checked;
  }

  setFilter(event: any) {
    console.log(this.newFilters);
    const filterVal = Number(event.target.value);
    switch (this.selectedFilter) {
      case this.filterNames.alchprof:
        this.filterValues.alchprof.filter = filterVal;
        break;
      case this.filterNames.margin:
        this.filterValues.margin.filter = filterVal;
        break;
      case this.filterNames.buyLimit:
        this.filterValues.buyLimit.filter = filterVal;
        break;
      case this.filterNames.highVolume:
        this.filterValues.highVolume.filter = filterVal;
        break;
      case this.filterNames.lowVolume:
        this.filterValues.lowVolume.filter = filterVal;
        break;

      default:
        console.error(
          'Error saving filter. Filter field name does not exist...'
        );
        break;
    }
  }

  getMaxValue(selectedFilter: string): number {
    switch (selectedFilter) {
      case this.filterNames.alchprof:
        return this.filterValues.alchprof.max;
      case this.filterNames.margin:
        return this.filterValues.margin.max;
      case this.filterNames.buyLimit:
        return this.filterValues.buyLimit.max;
      case this.filterNames.highVolume:
        return this.filterValues.highVolume.max;
      case this.filterNames.lowVolume:
        return this.filterValues.lowVolume.max;
      default:
        return 100;
    }
  }

  clearFilters(): void {
    this.filterValues.alchprof.filter = 1;
    this.filterValues.buyLimit.filter = 1;
    this.filterValues.highVolume.filter = 1;
    this.filterValues.lowVolume.filter = 1;
    this.filterValues.margin.filter = 5;
    this.updateFilter.emit(this.filterValues);
    this.selectedFilter = null;
  }

  onFilter(): void {
    this.updateFilter.emit(this.filterValues);
    this.selectedFilter = null;
  }
}
