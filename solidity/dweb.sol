// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./Idweb.sol";

contract DWeb is IDWeb, Ownable{
    string private _identity;
    string private _boostrap;
    bool private inited;

    constructor() {
        inited = false;
    }

    function identity() external view returns (string memory) {
        return _identity;
    }

    function boostrap() external view returns (string memory) {
        return _boostrap;
    }

    function initial(string calldata identity) external onlyOwner returns (bool) {
        require(!inited, "DWeb: contract has been inited");

        _identity = identity;
        inited = true;
        return true;
    }

    function update(string calldata newIdent) external onlyOwner returns (bool) {
        require(inited, "DWeb: contract require to init");

        _identity = newIdent;
        return true;
    }

    function join(string calldata link) external onlyOwner returns (bool) {
        _boostrap = link;
        return true;
    }
}