import { Component, Input, OnInit, EventEmitter, Output } from '@angular/core';
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
@Component({
  selector: 'app-listing-filters',
  standalone: true,
  imports: [
    MatButtonToggleModule,
    MatSelectModule,
    MatCheckboxModule,
    MatInputModule,
  ],
  templateUrl: './listing-filters.component.html',
  styleUrl: './listing-filters.component.scss',
})
export class ListingFiltersComponent {
  @Input() filterValues!: IFilters;
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
    let ret = 0;
    switch (selectedFilter) {
      case this.filterNames.alchprof:
        ret = this.filterValues.alchprof.max;
        break;
      case this.filterNames.margin:
        ret = this.filterValues.margin.max;
        break;
      case this.filterNames.buyLimit:
        ret = this.filterValues.buyLimit.max;
        break;
      case this.filterNames.highVolume:
        ret = this.filterValues.highVolume.max;
        break;
      case this.filterNames.lowVolume:
        ret = this.filterValues.lowVolume.max;
        break;
      default:
        ret = 100;
        break;
    }
    return Math.min(ret, 9999);
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
