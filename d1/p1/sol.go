package main

func main() {
	/*
	 * Strategy:
	 *   The naive method is to parse systematically through every pair until a hit is found - there's only one, anyway
	 *   Slightly less naive is to:
	 *     1. Sort low to high
	 *     2. Evaluate the sum of index 0 and index length-1
	 *     3. If the sum is > 2020, then one of the two numbers must be lowered - and the only way to accomplish this is to decrease the higher end; -- the higher end index
	 *     4. If the sum is > 2020, then one of the two numbers must be raised - and the only way to accomplish this is to increase the lower end; ++ the lower end index
	 */
}
