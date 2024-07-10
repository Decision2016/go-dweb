// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "./ownable.sol";

interface IDWeb{
    function identity() external view returns (string memory);
    function boostrap() external view returns (string memory);

    function initial(string calldata identity) external returns (bool);
    function update(string calldata identity) external returns (bool);
    function join(string calldata link) external returns (bool);
}