import { Component, Input, OnInit, EventEmitter, Output } from '@angular/core';
import {
  MatButtonToggleChange,
  MatButtonToggleModule,
} from '@angular/material/button-toggle';
import { MatSliderModule } from '@angular/material/slider';
import { MatSelectModule } from '@angular/material/select';
import {
  MatCheckboxChange,
  MatCheckboxModule,
} from '@angular/material/checkbox';
import { IFilters } from '../../interfaces/filters.interface';
@Component({
  selector: 'app-listing-filters',
  standalone: true,
  imports: [
    MatButtonToggleModule,
    MatSliderModule,
    MatSelectModule,
    MatCheckboxModule,
  ],
  templateUrl: './listing-filters.component.html',
  styleUrl: './listing-filters.component.scss',
})
export class ListingFiltersComponent implements OnInit {
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

  ngOnInit(): void {
    this.newFilters = {
      dataType: this.filterValues.dataType,
      alchprof: { max: 0, filter: 1 },
      margin: { max: 0, filter: 5 },
      buyLimit: { max: 0, filter: 1 },
      highVolume: { max: 0, filter: 1 },
      lowVolume: { max: 0, filter: 1 },
      members: true,
    };

    this.resetFilters = {
      dataType: this.filterValues.dataType,
      alchprof: { max: 0, filter: 1 },
      margin: { max: 0, filter: 5 },
      buyLimit: { max: 0, filter: 1 },
      highVolume: { max: 0, filter: 1 },
      lowVolume: { max: 0, filter: 1 },
      members: true,
    };
  }

  toggleView(event: MatButtonToggleChange) {
    this.newFilters.dataType = event.value;
  }

  toggleMembers(event: MatCheckboxChange) {
    this.newFilters.members = !event.checked;
  }

  setFilter(event: any) {
    const filterVal = Number(event.target.value);
    switch (this.selectedFilter) {
      case this.filterNames.alchprof:
        this.newFilters.alchprof.filter = filterVal;
        break;
      case this.filterNames.margin:
        this.newFilters.margin.filter = filterVal;
        break;
      case this.filterNames.buyLimit:
        this.newFilters.buyLimit.filter = filterVal;
        break;
      case this.filterNames.highVolume:
        this.newFilters.highVolume.filter = filterVal;
        break;
      case this.filterNames.lowVolume:
        this.newFilters.lowVolume.filter = filterVal;
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
    this.updateFilter.emit(this.resetFilters);
  }

  onFilter(): void {
    this.updateFilter.emit(this.newFilters);
  }
}
